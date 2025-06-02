'use client';

import { Accordion, AccordionItem } from '@heroui/react';
import '../styles/SidebarSection.css';

export default function SidebarSection({ title, children, defaultOpen = true, id }) {
  // Generate a fallback key/id if not provided
  const sectionId = id || title?.toLowerCase().replace(/\s+/g, '-');

  return (
    <section className="sidebar-section">
      <Accordion
        variant="shadow"
        defaultExpandedKeys={defaultOpen ? [sectionId] : []}
        className="sidebar-accordion"
      >
        <AccordionItem
          key={sectionId}
          aria-label={title}
          title={title}
          className="sidebar-section-title"
        >
          {children}
        </AccordionItem>
      </Accordion>
    </section>
  );
}
