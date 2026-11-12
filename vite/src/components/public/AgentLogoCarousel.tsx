import { useAgentLogos } from '@/hooks/useAgentLogos';
import { getStorageUrl } from '@/lib/utils';

export function AgentLogoCarousel() {
  const { logos, isLoading } = useAgentLogos(30);

  if (isLoading || logos.length === 0) return null;

  // ทำซ้ำ logos เพื่อให้ marquee วนต่อเนื่อง
  const duplicated = [...logos, ...logos];

  return (
    <div className="space-y-6">
      <div className="text-center space-y-2">
        <h2 className="text-2xl md:text-3xl font-bold">ตัวแทนจำหน่ายที่ไว้วางใจเรา</h2>
        <p className="text-muted-foreground">บริษัททัวร์ชั้นนำที่ร่วมเป็นพันธมิตรกับเรา</p>
      </div>
      <div className="relative overflow-hidden">
        {/* Fade edges */}
        <div className="pointer-events-none absolute inset-y-0 left-0 w-16 bg-gradient-to-r from-background to-transparent z-10" />
        <div className="pointer-events-none absolute inset-y-0 right-0 w-16 bg-gradient-to-l from-background to-transparent z-10" />

        <div className="flex items-center gap-8 animate-marquee hover:[animation-play-state:paused]">
          {duplicated.map((logo, i) => (
            <img
              key={i}
              src={getStorageUrl(logo)}
              alt="Agent logo"
              className="h-12 w-auto max-w-[120px] object-contain flex-shrink-0"
              onError={(e) => { (e.currentTarget as HTMLImageElement).style.display = 'none' }}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
