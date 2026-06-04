import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  ArrowRight, Shield, Zap, Lock,
  CheckCircle2, TrendingUp, Database, Cpu,
  ChevronDown, ChevronUp
} from 'lucide-react';
import { useState } from 'react';

const stats = [
  { value: '1.2M+', label: 'Sertifika Sayısı' },
  { value: '500+', label: 'Kurum Sayısı' },
  { value: '120', label: 'Ülke' },
];

const features = [
  {
    icon: Shield,
    title: 'Sertifika Güvenliği',
    desc: 'Diplomalarınız blockchain üzerinde değiştirilemez şekilde saklanır.',
  },
  {
    icon: Zap,
    title: 'Zaman Kazanın',
    desc: 'Dakikalar içinde dijital diploma oluşturun ve paylaşın.',
  },
  {
    icon: TrendingUp,
    title: 'Düşük Maliyet',
    desc: 'Geleneksel yöntemlere kıyasla %90 daha düşük maliyet.',
  },
  {
    icon: Lock,
    title: 'Kayıt Güvencesi',
    desc: 'Arweave ile kalıcı kayıt; veri hiçbir zaman kaybolmaz.',
  },
];

const steps = [
  { num: '01', title: 'Veri Yükleme', desc: 'Öğrenci bilgilerini ve diploma belgesini sisteme yükleyin.' },
  { num: '02', title: 'Şifreleme', desc: 'AES-256 ile belge şifrelenerek Arweave\'e yüklenir.' },
  { num: '03', title: 'Kayıt', desc: 'İşlem hash\'i Polygon ağına yazılır ve doğrulanır.' },
  { num: '04', title: 'Doğrulama', desc: 'TX numarası ile dünyanın her yerinden doğrulanabilir.' },
];

const faqs = [
  { q: 'Güvenlik için ne hizmetler sunulur?', a: 'Arweave kalıcı depolama ve Polygon blockchain ağını kullanıyoruz. Tüm veriler AES-256 şifreleme ile korunmaktadır.' },
  { q: 'Blockchain diplomalarına nasıl güvenir güvenilir?', a: 'Merkezi bir otorite olmadan doğrulanabilir; blok zinciri üzerindeki kayıtlar değiştirilemez ve herkese açıktır.' },
  { q: 'Mevcut sistemlere entegre edilebilir mi?', a: 'Evet, REST API aracılığıyla mevcut üniversite yönetim sistemlerine kolayca entegre edilebilir.' },
];

export default function Home() {
  const [openFaq, setOpenFaq] = useState<number | null>(null);

  return (
    <div className="flex flex-col w-full">
      {/* ── Hero (dark) ── */}
      <section className="min-h-screen bg-primary-container text-white pt-16 flex flex-col">
        <div className="flex-1 flex flex-col items-center justify-center text-center px-6 py-24">
          <motion.div
            initial={{ opacity: 0, y: 24 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="max-w-3xl mx-auto"
          >
            <span className="inline-block px-4 py-1.5 mb-8 text-xs font-semibold tracking-widest uppercase rounded-full bg-secondary-container/20 text-secondary-container border border-secondary-container/30">
              Blockchain Diploma Doğrulama
            </span>

            <h1 className="text-5xl md:text-6xl font-bold leading-tight mb-6">
              Akademik Geleceğiniz,{' '}
              <span className="text-secondary-container">Blockchain</span>
              {' '}ile Güvende
            </h1>

            <p className="text-lg text-white/60 max-w-xl mx-auto mb-10 leading-relaxed">
              Diplomalarınızı dijital ortamda güvenle saklayın, akademik geçmişinizi
              dünya genelinde kolayca doğrulatın.
            </p>

            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
              <Link
                to="/giris"
                className="group flex items-center gap-2 px-8 py-3.5 bg-white text-primary rounded-lg font-semibold hover:bg-white/90 transition-all"
              >
                Hemen Başla
                <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </Link>
              <a
                href="#how"
                className="px-8 py-3.5 border border-white/20 text-white rounded-lg font-semibold hover:bg-white/10 transition-all"
              >
                Nasıl Çalışır?
              </a>
            </div>

            <div className="flex flex-col sm:flex-row items-center justify-center gap-6 text-sm text-white/50">
              {['Arweave Kalıcı Depolama', 'Polygon Ağı Doğrulaması', 'AES-256 Şifreleme'].map((t) => (
                <span key={t} className="flex items-center gap-1.5">
                  <CheckCircle2 className="w-3.5 h-3.5 text-secondary-container" />
                  {t}
                </span>
              ))}
            </div>
          </motion.div>
        </div>

        {/* Stats */}
        <div className="border-t border-white/10 bg-white/5">
          <div className="max-w-[1280px] mx-auto px-6 py-10 grid grid-cols-3 gap-8 text-center">
            {stats.map(({ value, label }) => (
              <div key={label}>
                <p className="text-3xl font-bold text-white mb-1">{value}</p>
                <p className="text-sm text-white/50">{label}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Neden Biz ── */}
      <section id="about" className="bg-white py-24 px-6">
        <div className="max-w-[1280px] mx-auto">
          <div className="text-center mb-16">
            <span className="inline-block px-3 py-1 mb-4 text-xs font-semibold tracking-widest uppercase rounded-full bg-secondary/10 text-secondary border border-secondary/20">
              NEDEN BİZ?
            </span>
            <h2 className="text-4xl font-bold text-on-surface mb-4">Neden Biz?</h2>
            <p className="text-on-surface-variant max-w-lg mx-auto">
              Geleneksel diploma süreçlerini dijital, güvenli ve anında doğrulanabilir hale getiriyoruz.
            </p>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {features.map(({ icon: Icon, title, desc }) => (
              <motion.div
                key={title}
                whileHover={{ y: -4 }}
                className="p-6 rounded-2xl border border-outline-variant/40 bg-surface hover:border-on-tertiary-container/30 transition-all card-shadow"
              >
                <div className="w-10 h-10 rounded-xl bg-surface-container flex items-center justify-center mb-4">
                  <Icon className="w-5 h-5 text-on-tertiary-container" />
                </div>
                <h3 className="font-semibold text-on-surface mb-2">{title}</h3>
                <p className="text-sm text-on-surface-variant leading-relaxed">{desc}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Gerçek Bağlantılar ── */}
      <section className="bg-surface py-24 px-6">
        <div className="max-w-[1280px] mx-auto">
          <h2 className="text-4xl font-bold text-on-surface text-center mb-4">
            Gerçek Bağlantılar, Güvenilir Sertifikalar
          </h2>
          <p className="text-on-surface-variant text-center mb-16 max-w-lg mx-auto">
            Sahte diploma geçmişine gerçek bir son vermek için.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="p-8 rounded-2xl bg-white border border-outline-variant/40 card-shadow">
              <div className="w-10 h-10 rounded-xl bg-tertiary-container flex items-center justify-center mb-5">
                <Database className="w-5 h-5 text-on-tertiary-container" />
              </div>
              <h3 className="text-xl font-semibold text-on-surface mb-3">Değiştirilemez Doğrulama</h3>
              <p className="text-on-surface-variant leading-relaxed">
                Blok zincirindeki her kayıt şifrelenmiş ve birbirine bağlıdır. Hiçbir kayıt geri alınamaz,
                değiştirilemez veya silinemeez.
              </p>
              <Link to="/dogrula" className="inline-flex items-center gap-1 mt-5 text-sm text-on-tertiary-container font-medium hover:gap-2 transition-all">
                Şimdi Dene <ArrowRight className="w-4 h-4" />
              </Link>
            </div>
            <div className="p-8 rounded-2xl bg-white border border-outline-variant/40 card-shadow">
              <div className="w-10 h-10 rounded-xl bg-secondary/10 flex items-center justify-center mb-5">
                <Cpu className="w-5 h-5 text-secondary" />
              </div>
              <h3 className="text-xl font-semibold text-on-surface mb-3">Akıllı Aşılama Filtresi</h3>
              <p className="text-on-surface-variant leading-relaxed">
                Yapay zeka destekli sahte belge tespit sistemi, yükleme anında belgeyi doğrular
                ve potansiyel manipülasyonları tespit eder.
              </p>
              <Link to="/giris" className="inline-flex items-center gap-1 mt-5 text-sm text-on-tertiary-container font-medium hover:gap-2 transition-all">
                Şimdi Dene <ArrowRight className="w-4 h-4" />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* ── 4 Adım ── */}
      <section id="how" className="bg-white py-24 px-6">
        <div className="max-w-[1280px] mx-auto">
          <h2 className="text-4xl font-bold text-on-surface text-center mb-4">
            Toplam Güvenlik İçin 4 Adım
          </h2>
          <p className="text-on-surface-variant text-center mb-16 max-w-md mx-auto">
            Kolay kurulum, güçlü altyapı, küresel erişim ile başlayabilirsiniz.
          </p>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {steps.map(({ num, title, desc }) => (
              <div key={num} className="relative p-6 rounded-2xl bg-surface border border-outline-variant/40">
                <span className="text-4xl font-bold text-outline-variant/40 mb-4 block">{num}</span>
                <h3 className="font-semibold text-on-surface mb-2">{title}</h3>
                <p className="text-sm text-on-surface-variant leading-relaxed">{desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── FAQ ── */}
      <section className="bg-surface py-24 px-6">
        <div className="max-w-[768px] mx-auto">
          <h2 className="text-4xl font-bold text-on-surface text-center mb-4">Akılınıza Takılan</h2>
          <p className="text-on-surface-variant text-center mb-12">Sık sorulan sorular.</p>
          <div className="space-y-3">
            {faqs.map(({ q, a }, i) => (
              <div key={i} className="rounded-xl border border-outline-variant/40 bg-white overflow-hidden">
                <button
                  onClick={() => setOpenFaq(openFaq === i ? null : i)}
                  className="w-full flex items-center justify-between px-6 py-4 text-left"
                >
                  <span className="font-medium text-on-surface text-sm">{q}</span>
                  {openFaq === i
                    ? <ChevronUp className="w-4 h-4 text-outline shrink-0" />
                    : <ChevronDown className="w-4 h-4 text-outline shrink-0" />}
                </button>
                {openFaq === i && (
                  <div className="px-6 pb-4 text-sm text-on-surface-variant leading-relaxed">{a}</div>
                )}
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── CTA ── */}
      <section className="bg-primary-container py-24 px-6 text-center text-white">
        <div className="max-w-[768px] mx-auto">
          <h2 className="text-4xl font-bold mb-4">
            Geleceği Güvence Altına Almaya Hazır Mısınız?
          </h2>
          <p className="text-white/60 mb-10">
            500'den fazla kurumun emiyor olduğu doğrulanabilir sertifikalar ile akademik güvenceye bugün başlayın.
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              to="/giris"
              className="px-8 py-3.5 bg-white text-primary rounded-lg font-semibold hover:bg-white/90 transition-all"
            >
              Hemen Başla
            </Link>
            <Link
              to="/dogrula"
              className="px-8 py-3.5 border border-white/20 rounded-lg font-semibold hover:bg-white/10 transition-all"
            >
              Demo Dene
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
}
