import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import {
  FileText,
  Image as ImageIcon,
  MapPin,
} from 'lucide-react';
import { PublicPeriodsAccordion } from './PublicPeriodsAccordion';
import { getStorageUrl } from '@/lib/utils';
import type { Tour } from '@/types/tours-regular';

interface PublicTourCardProps {
  tour: Tour;
}

export function PublicTourCard({ tour }: PublicTourCardProps) {
  const handleDownload = (url: string, filename: string) => {
    // ส่ง filename ผ่าน query param เพื่อให้ backend ตั้ง Content-Disposition header
    const separator = url.includes('?') ? '&' : '?';
    const downloadUrl = `${url}${separator}filename=${encodeURIComponent(filename)}`;

    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <Card className="hover:shadow-md transition-shadow p-0">
      <Accordion type="single" collapsible className="w-full">
        <AccordionItem value={tour.id} className="border-0">
          <CardContent className="p-4 pb-0">
            <div className="grid grid-cols-1 md:grid-cols-[200px_1fr] gap-4">
              {/* Thumbnail */}
              <div className="space-y-2">
                {tour.thumbnail ? (
                  <img
                    src={getStorageUrl(tour.thumbnail)}
                    alt={tour.title}
                    className="w-full h-full md:h-full object-cover rounded-md"
                  />
                ) : (
                  <div className="w-full h-32 md:h-full bg-muted rounded-md flex items-center justify-center">
                    <ImageIcon className="h-8 w-8 text-muted-foreground" />
                  </div>
                )}
              </div>

              {/* Tour Info */}
              <div className="space-y-4">
                {/* Header */}
                <div className="space-y-2">
                  <div className="flex items-start justify-between">
                    <div>
                      <h3 className="text-lg font-semibold leading-tight">{tour.title}</h3>
                      <p className="text-sm text-muted-foreground">รหัสทัวร์: {tour.code}</p>
                    </div>
                  </div>

                  {tour.city_title && (
                    <div className="flex items-center text-sm text-muted-foreground">
                      <MapPin className="h-4 w-4 mr-1" />
                      {tour.city_title}
                    </div>
                  )}
                </div>

                {/* Detail */}
                {tour.detail && (
                  <p className="text-sm text-muted-foreground line-clamp-3 leading-relaxed">
                    {tour.detail}
                  </p>
                )}

                {/* Files */}
                <div className="flex flex-wrap gap-2">
                  {tour.file_pdf && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDownload(getStorageUrl(tour.file_pdf!), `${tour.code}_brochure.pdf`)}
                      className="bg-muted text-foreground text-md font-bold cursor-pointer"
                    >
                      <FileText className="h-3 w-3 mr-1 text-primary" />
                      PDF
                    </Button>
                  )}

                  {tour.file_word && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDownload(getStorageUrl(tour.file_word!), `${tour.code}_document.docx`)}
                      className="bg-muted text-foreground text-md font-bold cursor-pointer"
                    >
                      <FileText className="h-3 w-3 mr-1 text-chart-1" />
                      DOCX
                    </Button>
                  )}

                  {tour.thumbnail && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDownload(getStorageUrl(tour.thumbnail!), `${tour.code}_banner.jpg`)}
                      className="bg-muted text-foreground text-md font-bold cursor-pointer"
                    >
                      <ImageIcon className="h-3 w-3 mr-1" />
                      รูปภาพ
                    </Button>
                  )}
                </div>
              </div>
            </div>
          </CardContent>

          {/* Accordion Trigger */}
          <AccordionTrigger className="px-4 py-2 m-4 border-2 cursor-pointer bg-muted items-center hover:no-underline">
            <div className="flex items-center space-x-4 text-sm ">
              <div className="leading-none">
                <span className="block text-xm  text-foreground">Periods </span>
                <span className="block text-xl font-bold text-foreground leading-none">{tour.available_periods.length} รายการ</span>
              </div>
              {tour.min_start_date && (
                <div className="flex items-center text-foreground">
                  <div className="leading-none">
                    <span className="block text-xm  text-foreground">เร็วที่สุด </span>
                    <span className="block text-xl font-bold text-foreground leading-none">{new Date(tour.min_start_date).toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: 'numeric' })}</span>
                  </div>
                </div>
              )}
            </div>
          </AccordionTrigger>

          {/* Accordion Content - ไม่แสดง commission */}
          <AccordionContent className="px-4 pb-4">
            <PublicPeriodsAccordion tour={tour} />
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </Card>
  );
}
