// components/Image.tsx
import React, { useState, useEffect } from 'react'

interface Props extends React.ImgHTMLAttributes<HTMLImageElement> {
  src: string
  alt: string
  fallbackSrc?: string // ใช้เมื่อโหลดภาพหลักไม่สำเร็จ
  disableCache?: boolean // เพิ่ม prop สำหรับการปิดใช้งานการแคช
}

export const Image: React.FC<Props> = ({
  src,
  alt,
  fallbackSrc,
  loading = 'lazy',
  className,
  style,
  disableCache = false, // ค่าเริ่มต้นคือไม่ปิดการแคช
  ...rest
}) => {
  // เพิ่ม cache buster หากต้องการปิดการแคช
  const cacheBusterSrc = disableCache ? `${src}?t=${Date.now()}` : src;
  
  const [imgSrc, setImgSrc] = useState(cacheBusterSrc);

  // สำคัญมาก: อัพเดต imgSrc เมื่อ src เปลี่ยน
  useEffect(() => {
    // ใช้ cache buster หากต้องการปิดการแคช
    setImgSrc(disableCache ? `${src}?t=${Date.now()}` : src);
  }, [src, disableCache]);

  return (
    <img
      src={imgSrc}
      alt={alt}
      loading={loading}
      className={className}
      style={style}
      onError={() => {
        if (fallbackSrc) setImgSrc(fallbackSrc)
      }}
      {...rest}
    />
  )
}