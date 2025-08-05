# Layouts

The system uses a hierarchical, three-tiered layout mechanism to render content.

1.  **Content-Specific Layout (Highest Priority):** Individual content items can specify a particular layout to use. This provides the finest-grained control, allowing a post to have a unique design, overriding any section-level settings.

2.  **Section-Assigned Layout:** Every piece of content belongs to a section, which defines its URL path (e.g., `/documentation`). The site's root is also a section. If a content item doesn't have a specific layout assigned, it uses the one associated with its section. A single layout can be reused across multiple sections, and it can be the default one loaded from the seeder.

3.  **Embedded Fallback Layout (Lowest Priority):** If the layout specified by either the content or the section is not found in the database (e.g., it was deleted), the renderer defaults to using an embedded layout template embedded in the application binary (`assets/template/layout/layout.tmpl`). This serves as a failsafe to ensure that a view can always be rendered.

An editable, database-persisted version of the default layout is initially created via a data seed (`assets/seed/sqlite/20250707102435-ssg-add-core-data.json`), providing a ready-to-use, customizable template for sections.
