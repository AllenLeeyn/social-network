import { useState, useEffect } from "react";
import Image from "next/image";

export default function DynamicImage({ src, alt }) {
  const [dimensions, setDimensions] = useState(null);

  useEffect(() => {
    const img = new window.Image();
    img.src = src;
    img.onload = () => {
      setDimensions({
        width: img.width,
        height: img.height,
      });
    };
  }, [src]);

  if (!dimensions) return <div style={{ height: 200 }}>Loading image...</div>;

  return (
    <Image
      src={src}
      alt={alt}
      width={dimensions.width}
      height={dimensions.height}
      style={{ maxWidth: "50%", height: "auto" }}
    />
  );
}
