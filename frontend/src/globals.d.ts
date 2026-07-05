/// <reference types="astro/client" />

declare namespace App {
  interface Locals {
    blogContext?: {
      id: number;
      title: string;
      category?: string | null;
      searchTags: string[];
      createdAt: string;
      contentMd: string;
      contentHtml: string;
      authorId: number;
      authorName: string;
      authorDisplayName?: string | null;
    };
  }
}
