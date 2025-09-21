// Template loader utility for GenerateCode component

import { clientTemplates, serverTemplates } from "./templateRegistry";

interface TemplateCache {
  [key: string]: string;
}

class TemplateLoader {
  private clientTemplatesCache: TemplateCache = {};
  private serverTemplatesCache: TemplateCache = {};

  // Load a client template
  async loadClientTemplate(platform: string): Promise<string> {
    const cacheKey = `client-${platform}`;

    if (this.clientTemplatesCache[cacheKey]) {
      return this.clientTemplatesCache[cacheKey];
    }

    try {
      const template = clientTemplates[platform];
      if (!template) {
        throw new Error(`Client template not found for platform: ${platform}`);
      }

      this.clientTemplatesCache[cacheKey] = template;
      return template;
    } catch (error) {
      console.error(`Error loading client template ${platform}:`, error);
      return "// Template not found or failed to load";
    }
  }

  // Load a server template
  async loadServerTemplate(language: string): Promise<string> {
    const cacheKey = `server-${language}`;

    if (this.serverTemplatesCache[cacheKey]) {
      return this.serverTemplatesCache[cacheKey];
    }

    try {
      const template = serverTemplates[language];
      if (!template) {
        throw new Error(`Server template not found for language: ${language}`);
      }

      this.serverTemplatesCache[cacheKey] = template;
      return template;
    } catch (error) {
      console.error(`Error loading server template ${language}:`, error);
      return "// Template not found or failed to load";
    }
  }

  // Process template with variables
  processTemplate(
    template: string,
    variables: { [key: string]: string }
  ): string {
    let processed = template;

    for (const [key, value] of Object.entries(variables)) {
      const regex = new RegExp(`{{${key}}}`, "g");
      processed = processed.replace(regex, value);
    }

    return processed;
  }

  // Clear cache (useful for development)
  clearCache(): void {
    this.clientTemplatesCache = {};
    this.serverTemplatesCache = {};
  }
}

// Export singleton instance
export const templateLoader = new TemplateLoader();
