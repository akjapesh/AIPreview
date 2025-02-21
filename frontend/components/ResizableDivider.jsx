"use client";

import React, { useCallback, useRef } from "react";
import { useStyletron } from "baseui";

export const ResizableDivider = ({ onResize }) => {
  const [css] = useStyletron();
  const dividerRef = useRef(null);

  const handleMouseDown = useCallback(
    (e) => {
      e.preventDefault();

      const startX = e.clientX;
      const startWidth =
        dividerRef.current?.previousElementSibling?.getBoundingClientRect()
          .width || 0;

      const handleMouseMove = (moveEvent) => {
        const delta = moveEvent.clientX - startX;
        const newWidth = startWidth + delta;
        const containerWidth =
          dividerRef.current?.parentElement?.getBoundingClientRect().width || 0;

        // Calculate percentage width and clamp between 20% and 80%
        const percentage = Math.min(
          Math.max((newWidth / containerWidth) * 100, 20),
          80
        );
        onResize(`${percentage}%`);
      };

      const handleMouseUp = () => {
        document.removeEventListener("mousemove", handleMouseMove);
        document.removeEventListener("mouseup", handleMouseUp);
      };

      document.addEventListener("mousemove", handleMouseMove);
      document.addEventListener("mouseup", handleMouseUp);
    },
    [onResize]
  );

  return (
    <div
      ref={dividerRef}
      className={css({
        width: "8px",
        cursor: "col-resize",
        height: "100%",
        backgroundColor: "rgba(0, 0, 0, 0.1)",
        ":hover": {
          backgroundColor: "rgba(0, 0, 0, 0.2)",
        },
      })}
      onMouseDown={handleMouseDown}
    />
  );
};
