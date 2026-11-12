import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Flame,
  Plane,
  Search,
  LogIn,
  UserPlus,
} from 'lucide-react';
import { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { PublicToursRegularList } from '@/components/public/PublicToursRegularList';
import { PublicTourCard } from '@/components/public/PublicTourCard';
import { CityFilter } from '@/components/tours-regular/CityFilter';
import { AgentLogoCarousel } from '@/components/public/AgentLogoCarousel';
import { PublicBannerCarousel } from '@/components/public/PublicBannerCarousel';
import { usePublicBanner } from '@/hooks/usePublicBanner';
import { useToursHot } from '@/hooks/useToursHot';
import { useAuth } from '@/hooks/useAuth';
import { ModeToggle } from '@/theme/mode-toggle';
export default function HomePage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchText, setSearchText] = useState(searchParams.get('search') || '');
  const [debouncedSearch, setDebouncedSearch] = useState(searchText);
  const [selectedCityId, setSelectedCityId] = useState<string | undefined>();
  const navigate = useNavigate();
  const { isAuthenticated, user } = useAuth();

  // ดึง banner ที่ active (public - ไม่ต้อง auth) เฉพาะ position=public
  const { banners: bannersList } = usePublicBanner({ is_active: true, limit: 10, position: 'public' });

  // ดึงทัวร์ไฟไหม้ (public route - ไม่ต้อง auth)
  const { tours: hotTours, isLoading: hotToursLoading } = useToursHot(
    { status: 'publish', available_only: true, limit: 6 },
    { enabled: true }
  );

  // Debounce search - 500ms delay
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(searchText);
    }, 500);

    return () => clearTimeout(timer);
  }, [searchText]);

  // Update URL when debounced search changes
  useEffect(() => {
    if (debouncedSearch.trim()) {
      setSearchParams({ search: debouncedSearch.trim() });
    } else {
      // ลบ search param เมื่อค้นหาว่าง
      setSearchParams({});
    }
  }, [debouncedSearch, setSearchParams]);

  const handleDashboard = () => {
    if (user?.role === 'admin' || user?.role === 'super') {
      navigate('/admin/dashboard');
    } else if (user?.role === 'agent') {
      navigate('/agent/dashboard');
    } else if (user?.role === 'sales') {
      navigate('/sale/dashboard');
    } else {
      navigate('/dashboard');
    }
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-16 items-center justify-between px-4 mx-auto max-w-6xl">
          {/* Logo */}
          <Link to="/" className="flex items-center space-x-1.5">
            <img src="/images/id/logo.png" alt="TIAOJIN" className="h-12 w-auto" />
            <div className="flex flex-col">
              <span className="font-bold text-xl leading-tight">TIAOJIN</span>
              <span className="text-xs text-muted-foreground leading-tight">ทัวร์จีนต้องเที่ยวจีน</span>
            </div>
          </Link>

          {/* Navigation */}
          <nav className="flex items-center gap-2">
            <ModeToggle />
            {isAuthenticated ? (
              <Button onClick={handleDashboard} variant="default" size="sm">
                แดชบอร์ด
              </Button>
            ) : (
              <Button asChild size="sm">
                <Link to="/login">
                  เข้าสู่ระบบ
                </Link>
              </Button>
            )}
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <section className="relative pt-4 pb-8 md:pt-6 md:pb-12 bg-gradient-to-b from-primary/5 to-background">
        <div className="container px-4 mx-auto max-w-6xl">
          <div className="text-center space-y-6">
            {/* Banner Carousel */}
            {bannersList && bannersList.length > 0 && (
              <PublicBannerCarousel banners={bannersList} />
            )}

            <h1 className="text-4xl md:text-5xl font-bold tracking-tight">
              ค้นหา<span className="text-primary">ทัวร์จีน</span>ที่ใช่สำหรับคุณ
            </h1>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              รวมโปรแกรมทัวร์จีนคุณภาพ พร้อมราคาพิเศษสำหรับตัวแทนจำหน่าย
            </p>

            {/* Search Input - Auto search with debounce */}
            <div className="max-w-xl mx-auto">
              <div className="relative">
                <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
                <Input
                  type="text"
                  placeholder="ค้นหาทัวร์จีน... เช่น ปักกิ่ง, เซี่ยงไฮ้, กุ้ยหลิน"
                  value={searchText}
                  onChange={(e) => setSearchText(e.target.value)}
                  className="pl-12 h-14 text-base rounded-full border-2 focus:border-primary"
                />
                {searchText && (
                  <button
                    type="button"
                    onClick={() => setSearchText('')}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  >
                    ✕
                  </button>
                )}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Hot Tours Section */}
      {!hotToursLoading && hotTours.length > 0 && (
        <section className="py-8">
          <div className="container px-4 mx-auto max-w-6xl">
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center gap-2">
                <Flame className="h-5 w-5 text-destructive" />
                <h2 className="text-xl font-semibold">โปรไฟไหม้</h2>
              </div>
              {isAuthenticated && (
                <Button variant="ghost" size="sm" asChild>
                  <Link to="/agent/tours-hot">ดูทั้งหมด →</Link>
                </Button>
              )}
            </div>
            <div className="grid grid-cols-1 gap-4">
              {hotTours.map((tour) => (
                <PublicTourCard key={tour.id} tour={tour} />
              ))}
            </div>
          </div>
        </section>
      )}

      {/* (Banner moved to top of content) */}

      {/* Tours Section */}
      <section className="py-12">
        <div className="container px-4 mx-auto max-w-6xl">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-2">
              <Plane className="h-5 w-5 text-primary" />
              <h2 className="text-xl font-semibold">โปรแกรมทัวร์จีนทั้งหมด</h2>
            </div>
          </div>

          {/* City Filter */}
          <div className="mb-6">
            <CityFilter
              selectedCityId={selectedCityId}
              onCityChange={setSelectedCityId}
            />
          </div>

          {/* Tours List - ไม่แสดง commission */}
          <PublicToursRegularList
            key={selectedCityId || 'all'}
            initialParams={{
              status: 'publish',
              available_only: true,
              ...(selectedCityId ? { city_id: selectedCityId } : {}),
            }}
          />
        </div>
      </section>

      {/* Agent Logo Carousel — แสดงตลอดทั้ง login และไม่ login */}
      <section className="py-10 bg-background">
        <div className="container px-4 mx-auto max-w-6xl">
          <AgentLogoCarousel />
        </div>
      </section>

      {/* CTA Section */}
      {!isAuthenticated && (
        <section className="py-16 bg-primary/5">
          <div className="container px-4 mx-auto max-w-6xl text-center">
            <h2 className="text-2xl md:text-3xl font-bold mb-4">
              เป็นตัวแทนจำหน่ายกับเรา
            </h2>
            <p className="text-muted-foreground mb-8 max-w-xl mx-auto">
              ลงทะเบียนเป็นตัวแทนจำหน่ายวันนี้ เพื่อรับสิทธิพิเศษในการจองทัวร์และค่าคอมมิชชั่น
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button size="lg" asChild>
                <Link to="/register">
                  <UserPlus className="h-5 w-5 mr-2" />
                  ลงทะเบียนตัวแทน
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild>
                <Link to="/login">
                  <LogIn className="h-5 w-5 mr-2" />
                  เข้าสู่ระบบ
                </Link>
              </Button>
            </div>
          </div>
        </section>
      )}

      {/* Footer */}
      <footer className="border-t bg-muted/50">
        <div className="container px-4 mx-auto max-w-6xl">
          {/* Main Footer */}
          <div className="py-12 grid grid-cols-1 md:grid-cols-3 gap-8">
            {/* Brand */}
            <div className="space-y-4">
              <div className="flex items-center space-x-2">
                <img src="/images/id/logo.png" alt="TIAOJIN" className="h-10 w-auto" />
                <div className="flex flex-col">
                  <span className="font-bold text-lg leading-tight">TIAOJIN</span>
                  <span className="text-xs text-muted-foreground">ทัวร์จีนต้องเที่ยวจีน</span>
                </div>
              </div>
              <p className="text-sm text-muted-foreground">
                รวมโปรแกรมทัวร์จีนคุณภาพ พร้อมราคาพิเศษสำหรับตัวแทนจำหน่าย
              </p>
            </div>

            {/* Quick Links */}
            <div className="space-y-4">
              <h3 className="font-semibold">บริการ</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li>ทัวร์ปักกิ่ง</li>
                <li>ทัวร์เซี่ยงไฮ้</li>
                <li>ทัวร์กุ้ยหลิน</li>
                <li>ทัวร์จางเจียเจี้ย</li>
              </ul>
            </div>

            {/* Contact */}
            <div className="space-y-4">
              <h3 className="font-semibold">ติดต่อเรา</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li>โทร: 02-XXX-XXXX</li>
                <li>Line: @tiaojin</li>
                <li>Email: info@tiaojin.com</li>
              </ul>
            </div>
          </div>

          {/* Copyright */}
          <div className="py-6 border-t border-border/50">
            <p className="text-center text-sm text-muted-foreground">
              © {new Date().getFullYear()} TIAOJIN. All rights reserved.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
