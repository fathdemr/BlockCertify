---
name: BlockCertify
colors:
  surface: '#f7f9fb'
  surface-dim: '#d8dadc'
  surface-bright: '#f7f9fb'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f2f4f6'
  surface-container: '#eceef0'
  surface-container-high: '#e6e8ea'
  surface-container-highest: '#e0e3e5'
  on-surface: '#191c1e'
  on-surface-variant: '#45464d'
  inverse-surface: '#2d3133'
  inverse-on-surface: '#eff1f3'
  outline: '#76777d'
  outline-variant: '#c6c6cd'
  surface-tint: '#565e74'
  primary: '#000000'
  on-primary: '#ffffff'
  primary-container: '#131b2e'
  on-primary-container: '#7c839b'
  inverse-primary: '#bec6e0'
  secondary: '#006c49'
  on-secondary: '#ffffff'
  secondary-container: '#6cf8bb'
  on-secondary-container: '#00714d'
  tertiary: '#000000'
  on-tertiary: '#ffffff'
  tertiary-container: '#001a42'
  on-tertiary-container: '#3980f4'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#dae2fd'
  primary-fixed-dim: '#bec6e0'
  on-primary-fixed: '#131b2e'
  on-primary-fixed-variant: '#3f465c'
  secondary-fixed: '#6ffbbe'
  secondary-fixed-dim: '#4edea3'
  on-secondary-fixed: '#002113'
  on-secondary-fixed-variant: '#005236'
  tertiary-fixed: '#d8e2ff'
  tertiary-fixed-dim: '#adc6ff'
  on-tertiary-fixed: '#001a42'
  on-tertiary-fixed-variant: '#004395'
  background: '#f7f9fb'
  on-background: '#191c1e'
  surface-variant: '#e0e3e5'
typography:
  display-lg:
    fontFamily: Inter
    fontSize: 48px
    fontWeight: '700'
    lineHeight: 56px
    letterSpacing: -0.02em
  display-lg-mobile:
    fontFamily: Inter
    fontSize: 36px
    fontWeight: '700'
    lineHeight: 44px
    letterSpacing: -0.02em
  headline-md:
    fontFamily: Inter
    fontSize: 30px
    fontWeight: '600'
    lineHeight: 38px
    letterSpacing: -0.01em
  headline-sm:
    fontFamily: Inter
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
  body-lg:
    fontFamily: Inter
    fontSize: 18px
    fontWeight: '400'
    lineHeight: 28px
  body-md:
    fontFamily: Inter
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-sm:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-mono:
    fontFamily: JetBrains Mono
    fontSize: 13px
    fontWeight: '500'
    lineHeight: 16px
    letterSpacing: 0.05em
  button:
    fontFamily: Inter
    fontSize: 16px
    fontWeight: '600'
    lineHeight: 100%
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  base: 4px
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
  2xl: 48px
  3xl: 64px
  container-max: 1280px
  gutter: 24px
  margin-mobile: 16px
---

## Brand & Style
The design system for this platform focuses on **Trust, Verifiability, and Modernity**. As a SaaS blockchain credentialing tool, the interface must project the institutional authority of a diploma while maintaining the technical agility of Web3. 

The aesthetic is a hybrid of **Modern Corporate** and **Glassmorphism**. It utilizes spacious layouts, high-quality typography, and subtle translucency to suggest a "digital vault" atmosphere. The emotional response should be one of absolute security and professional clarity, ensuring that both the issuer and the recipient feel the weight and validity of the digital credentials. The system is fully optimized for Turkish typographic nuances (like dotted and dotless 'i') and layout expansion.

## Colors
The palette is rooted in deep, authoritative blues to establish trust, contrasted with a vibrant "Security Green" used exclusively for success states and verified actions.

- **Primary (Deep Navy):** Used for navigation, headings, and core structural elements. It represents the "Immutable" nature of the blockchain.
- **Secondary (Security Green):** Reserved for verification badges, "Issue" buttons, and positive status indicators.
- **Tertiary (Action Blue):** Used for secondary actions, text links, and interactive states to keep the UI from feeling overly heavy.
- **Neutral:** A range of cool grays and "Clean White" surfaces provide the necessary breathing room for complex data tables and certificate previews.

## Typography
This design system utilizes **Inter** for all primary communication due to its exceptional legibility and support for Turkish character sets (ş, ı, ğ, ç, ö, ü). **JetBrains Mono** is introduced as a secondary label font for hash addresses, transaction IDs, and metadata to emphasize the technical precision of the platform.

Hierarchy is strictly enforced through weight. Headings use semi-bold and bold weights to ground the page, while body text remains in regular weights for maximum comfort during long reading sessions.

## Layout & Spacing
The layout follows a **12-column fluid grid** for desktop and a **4-column grid** for mobile. We utilize a soft 8px rhythm to ensure vertical harmony.

- **Desktop (1280px+):** Centered container with 24px gutters. Large internal paddings (32px-48px) within cards to maintain a premium feel.
- **Tablet (768px - 1024px):** Gutters reduce to 16px. Sidebars may collapse into a drawer or a compact icon bar.
- **Mobile (<768px):** Single column flow with 16px horizontal margins. Buttons and touch targets scale to a minimum height of 48px.

Spaciousness is a functional requirement; data density is managed through progressive disclosure to prevent "dashboard fatigue."

## Elevation & Depth
Depth is created through a combination of **Tonal Layering** and **Glassmorphism**. 

1.  **Base Layer:** The canvas uses the Neutral hex (`#F8FAFC`).
2.  **Surface Layer:** Cards and containers use pure white (`#FFFFFF`) with a very soft, 10% opacity blue-tinted shadow.
3.  **Glass Layer:** Overlays, modals, and navigation headers use a semi-transparent white (e.g., `rgba(255, 255, 255, 0.7)`) with a `backdrop-filter: blur(12px)`. This creates a sense of "layered security" where the user can always see the context of the platform beneath the current action.
4.  **Interactive Elevation:** On hover, primary cards lift slightly (y-offset increases) and the border-color shifts to the Tertiary blue.

## Shapes
The design system adopts a **Rounded** profile (`0.5rem` or `8px` base) to balance professionalism with modern software approachability. 

- **Standard Buttons & Inputs:** 8px radius.
- **Cards & Modals:** 16px (rounded-lg) to create a soft, contained look for certificates.
- **Status Chips & Badges:** Full pill-shape (`9999px`) to distinguish them clearly from interactive buttons.
- **Verification Seals:** Circular or heavy rounded squares to evoke the feeling of a physical stamp.

## Components
- **Buttons:** 
    - *Primary:* Filled Deep Navy (`#0F172A`) with white text. 
    - *Success:* Filled Security Green (`#10B981`) for "Verify" or "Issue" actions.
    - *Ghost:* Transparent with Action Blue border for secondary navigation.
- **Input Fields:** Use a subtle 1px border (`#E2E8F0`). On focus, the border transitions to Action Blue with a soft 4px outer glow. Labels are always positioned above the input in `body-sm` bold.
- **Certificates (Cards):** Featured with a `1px` inner stroke and a high-blur shadow. They should include a "Verification Badge" in the top-right corner using the Security Green.
- **Verification Status (Turkish Context):**
    - "Doğrulandı" (Verified): Green pill with a checkmark icon.
    - "Beklemede" (Pending): Amber pill.
    - "Geçersiz" (Invalid): Red pill.
- **Data Tables:** High-contrast rows with `label-mono` used for all blockchain-related strings (wallet addresses, hashes).
- **Navigation:** A glassmorphic top-bar that remains fixed, ensuring the platform's "Identity" is always visible.