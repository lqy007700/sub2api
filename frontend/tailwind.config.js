/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 主色调 - 苹果雾银克制蓝 (Apple Mist)
        primary: {
          50: '#eef6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#3b9eff',
          500: '#0071e3',
          600: '#0062c4',
          700: '#004fa3',
          800: '#073f7d',
          900: '#0c3361',
          950: '#081f3d'
        },
        // 辅助色 - 中性雾灰
        accent: {
          50: '#f8f9fb',
          100: '#f1f3f7',
          200: '#e2e6ee',
          300: '#cbd2dc',
          400: '#8d99a8',
          500: '#526071',
          600: '#3c4756',
          700: '#2a3340',
          800: '#1c232d',
          900: '#111827',
          950: '#0b0c0f'
        },
        // 深色模式背景 - 纯正曜石炭黑
        dark: {
          50: '#f5f7fa',
          100: '#e7ebf0',
          200: '#cbd2dc',
          300: '#9aa3b2',
          400: '#6e788b',
          500: '#475063',
          600: '#2e3442',
          700: '#222631',
          800: '#161920',
          900: '#101217',
          950: '#0b0c0f'
        }
      },
      fontFamily: {
        sans: [
          '-apple-system',
          'BlinkMacSystemFont',
          '"SF Pro Display"',
          '"SF Pro Text"',
          'system-ui',
          'Segoe UI',
          'PingFang SC',
          'Hiragino Sans GB',
          'Microsoft YaHei',
          'sans-serif'
        ],
        mono: ['"SF Mono"', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace']
      },
      boxShadow: {
        glass: '0 1px 2px rgba(15, 23, 42, 0.04), 0 12px 32px -8px rgba(15, 23, 42, 0.08)',
        'glass-sm': '0 1px 2px rgba(15, 23, 42, 0.04), 0 4px 12px -2px rgba(15, 23, 42, 0.05)',
        glow: '0 0 16px rgba(0, 113, 227, 0.14)',
        'glow-lg': '0 0 28px rgba(0, 113, 227, 0.2)',
        card: '0 1px 2px rgba(15, 23, 42, 0.03), 0 8px 24px -8px rgba(15, 23, 42, 0.06)',
        'card-hover': '0 2px 4px rgba(15, 23, 42, 0.04), 0 14px 32px -10px rgba(15, 23, 42, 0.09)',
        pill: '0 1px 3px rgba(15, 23, 42, 0.08), 0 0 0 1px rgba(15, 23, 42, 0.04)',
        'inner-glow': 'inset 0 1px 0 rgba(255, 255, 255, 0.85)'
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-primary': 'linear-gradient(135deg, #0071e3 0%, #0062c4 100%)',
        'gradient-dark': 'linear-gradient(135deg, #161920 0%, #0b0c0f 100%)',
        'gradient-glass':
          'linear-gradient(135deg, rgba(255,255,255,0.12) 0%, rgba(255,255,255,0.04) 100%)',
        'mesh-gradient':
          'radial-gradient(920px 380px at 50% -120px, rgba(0, 113, 227, 0.075) 0px, transparent 70%)'
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-in-right': 'slideInRight 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shimmer: 'shimmer 2s linear infinite',
        glow: 'glow 2s ease-in-out infinite alternate'
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' }
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideInRight: {
          '0%': { opacity: '0', transform: 'translateX(20px)' },
          '100%': { opacity: '1', transform: 'translateX(0)' }
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' }
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' }
        },
        glow: {
          '0%': { boxShadow: '0 0 16px rgba(56, 87, 214, 0.14)' },
          '100%': { boxShadow: '0 0 24px rgba(56, 87, 214, 0.24)' }
        }
      },
      backdropBlur: {
        xs: '2px'
      },
      borderRadius: {
        '4xl': '2rem'
      }
    }
  },
  plugins: []
}
