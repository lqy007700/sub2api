//go:build unit

package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/Wei-Shaw/sub2api/migrations"
	"github.com/stretchr/testify/require"
)

// This DSN must target an isolated PostgreSQL cluster and permit CREATE DATABASE.
// Each case owns a separate database because legacy migrations refer to public.
func TestForkFunctionalSyncPostgresMigrations(t *testing.T) {
	dsn := os.Getenv("FORK_MIGRATION_TEST_DSN")
	if dsn == "" {
		t.Skip("requires isolated FORK_MIGRATION_TEST_DSN with CREATE DATABASE permission")
	}
	base, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = base.Close() })

	for _, upgrade := range []bool{false, true} {
		t.Run(fmt.Sprintf("upgrade=%v", upgrade), func(t *testing.T) {
			name := fmt.Sprintf("fork_sync_%d", time.Now().UnixNano())
			_, err := base.Exec("CREATE DATABASE " + name)
			require.NoError(t, err)
			db, err := sql.Open("postgres", dsn+" dbname="+name)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = db.Close()
				_, dropErr := base.Exec("DROP DATABASE " + name + " WITH (FORCE)")
				require.NoError(t, dropErr)
			})
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			if upgrade {
				baseline := fstest.MapFS{}
				names, err := fs.Glob(migrations.FS, "*.sql")
				require.NoError(t, err)
				for _, filename := range names {
					if filename > "247_group_hidden_rate_multiplier.sql" {
						continue
					}
					data, err := migrations.FS.ReadFile(filename)
					require.NoError(t, err)
					baseline[filename] = &fstest.MapFile{Data: data}
				}
				require.NoError(t, applyMigrationsFS(ctx, db, baseline))
			} else {
				require.NoError(t, ApplyMigrations(ctx, db))
			}

			var groupID, userID int64
			require.NoError(t, db.QueryRowContext(ctx, `INSERT INTO groups(name,platform,rate_multiplier,hidden_rate_multiplier,allow_image_generation)
				VALUES('hidden-image-group','openai',1.25,2,true) RETURNING id`).Scan(&groupID))
			require.NoError(t, db.QueryRowContext(ctx, `INSERT INTO users(email,password_hash) VALUES('fork-sync@example.test','test-only') RETURNING id`).Scan(&userID))
			const key = "sk-fork-sync-test"
			_, err = db.ExecContext(ctx, `INSERT INTO api_keys(user_id,key,name,group_id) VALUES($1,$2,'fork-sync',$3)`, userID, key, groupID)
			require.NoError(t, err)
			require.NoError(t, ApplyMigrations(ctx, db))
			require.NoError(t, ApplyMigrations(ctx, db), "migrations must replay without checksum errors")

			var hidden, public float64
			var countBilling bool
			require.NoError(t, db.QueryRowContext(ctx, `SELECT hidden_rate_multiplier,rate_multiplier,image_count_billing FROM groups WHERE id=$1`, groupID).
				Scan(&hidden, &public, &countBilling))
			require.Equal(t, 2.0, hidden)
			require.Equal(t, 1.25, public)
			require.False(t, countBilling, "existing and new groups default to the previous billing behavior")
			for _, filename := range []string{"247_group_hidden_rate_multiplier.sql", "247_scoped_admin_api_keys.sql", "248_group_image_count_billing.sql", "249_group_auth_cache_image_count_billing.sql"} {
				var checksum string
				require.NoError(t, db.QueryRowContext(ctx, "SELECT checksum FROM schema_migrations WHERE filename=$1", filename).Scan(&checksum))
				data, err := migrations.FS.ReadFile(filename)
				require.NoError(t, err)
				require.Equal(t, fmt.Sprintf("%x", sha256.Sum256([]byte(strings.TrimSpace(string(data))))), checksum)
			}

			// SQL-only changes must invalidate both new and existing snapshot fields.
			for _, assignment := range []string{"image_count_billing=true", "hidden_rate_multiplier=3", "profit_control_enabled=true"} {
				_, err = db.ExecContext(ctx, "DELETE FROM auth_cache_invalidation_outbox")
				require.NoError(t, err)
				_, err = db.ExecContext(ctx, "UPDATE groups SET "+assignment+" WHERE id=$1", groupID)
				require.NoError(t, err)
				var cacheKey string
				require.NoError(t, db.QueryRowContext(ctx, "SELECT cache_key FROM auth_cache_invalidation_outbox").Scan(&cacheKey))
				require.Equal(t, fmt.Sprintf("%x", sha256.Sum256([]byte(key))), cacheKey)
			}
			_, err = db.ExecContext(ctx, "DELETE FROM auth_cache_invalidation_outbox")
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, "UPDATE groups SET name=name || '-renamed', image_count_billing=image_count_billing WHERE id=$1", groupID)
			require.NoError(t, err)
			var invalidations int
			require.NoError(t, db.QueryRowContext(ctx, "SELECT COUNT(*) FROM auth_cache_invalidation_outbox").Scan(&invalidations))
			require.Zero(t, invalidations, "unrelated edits must not enqueue invalidations")

			// Exercise merged Ent field indexes and the auth repository projection.
			client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
			repo := newGroupRepositoryWithSQL(client, db)
			group, err := repo.GetByID(ctx, groupID)
			require.NoError(t, err)
			require.Equal(t, 3.0, *group.HiddenRateMultiplier)
			require.True(t, group.ImageCountBilling)
			group.HiddenRateMultiplier = nil
			group.ImageCountBilling = false
			require.NoError(t, repo.Update(ctx, group))
			auth, err := newAPIKeyRepositoryWithSQL(client, db).GetByKeyForAuth(ctx, key)
			require.NoError(t, err)
			require.NotNil(t, auth.Group)
			require.Nil(t, auth.Group.HiddenRateMultiplier)
			require.False(t, auth.Group.ImageCountBilling)
			require.Equal(t, service.PlatformOpenAI, auth.Group.Platform)
		})
	}
}
