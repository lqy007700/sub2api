-- Optional actual billing rate; NULL retains the public group rate.
ALTER TABLE groups
    ADD COLUMN IF NOT EXISTS hidden_rate_multiplier DECIMAL(10,4) NULL;

ALTER TABLE groups
    ADD CONSTRAINT groups_hidden_rate_multiplier_positive
    CHECK (hidden_rate_multiplier IS NULL OR
           (hidden_rate_multiplier > 0 AND hidden_rate_multiplier <= 999999.9999));

-- Preserve the public multiplier at request time independently of actual billing.
-- Historical rows intentionally remain NULL and display their actual rate_multiplier.
ALTER TABLE usage_logs
    ADD COLUMN IF NOT EXISTS display_rate_multiplier DECIMAL(10,4) NULL;

-- Batch settlement must retain the submitted public group/discount multiplier.
ALTER TABLE batch_image_jobs
    ADD COLUMN IF NOT EXISTS display_rate_multiplier DECIMAL(10,4) NULL;

COMMENT ON COLUMN batch_image_jobs.display_rate_multiplier IS
    'User-visible final batch multiplier captured at submission';

COMMENT ON COLUMN groups.hidden_rate_multiplier IS
    'Optional actual billing multiplier; NULL uses the public rate_multiplier';
COMMENT ON COLUMN usage_logs.display_rate_multiplier IS
    'User-visible multiplier snapshot; NULL falls back to historical rate_multiplier';

-- Hidden billing rates are part of the auth snapshot. Preserve the latest
-- durable invalidation checks and include this new input for direct SQL edits.
CREATE OR REPLACE FUNCTION enqueue_group_auth_cache_invalidation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    target_group_id BIGINT;
BEGIN
    target_group_id := OLD.id;
    IF TG_OP = 'UPDATE'
       AND OLD.status IS NOT DISTINCT FROM NEW.status
       AND OLD.is_exclusive IS NOT DISTINCT FROM NEW.is_exclusive
       AND OLD.allow_image_generation IS NOT DISTINCT FROM NEW.allow_image_generation
       AND OLD.platform IS NOT DISTINCT FROM NEW.platform
       AND OLD.subscription_type IS NOT DISTINCT FROM NEW.subscription_type
       AND OLD.rate_multiplier IS NOT DISTINCT FROM NEW.rate_multiplier
       AND OLD.hidden_rate_multiplier IS NOT DISTINCT FROM NEW.hidden_rate_multiplier
       AND OLD.peak_rate_enabled IS NOT DISTINCT FROM NEW.peak_rate_enabled
       AND OLD.peak_start IS NOT DISTINCT FROM NEW.peak_start
       AND OLD.peak_end IS NOT DISTINCT FROM NEW.peak_end
       AND OLD.peak_rate_multiplier IS NOT DISTINCT FROM NEW.peak_rate_multiplier
       AND OLD.profit_control_enabled IS NOT DISTINCT FROM NEW.profit_control_enabled
       AND OLD.profit_min_margin IS NOT DISTINCT FROM NEW.profit_min_margin
       AND OLD.profit_safety_buffer IS NOT DISTINCT FROM NEW.profit_safety_buffer
       AND OLD.deleted_at IS NOT DISTINCT FROM NEW.deleted_at THEN
        RETURN NEW;
    END IF;

    INSERT INTO auth_cache_invalidation_outbox (cache_key)
    SELECT encode(sha256(convert_to(k.key, 'UTF8')), 'hex')
    FROM api_keys AS k
    WHERE k.group_id = target_group_id
      AND k.deleted_at IS NULL
      AND k.key <> '';
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$;
