// This script is for converting Firefox bookmarks.html (might work with Chromium) into data.json
// Place bookmarks.html in this folder and run this command:
// deno run -A --unstable convertBookmarksHTMLtoJSON.ts

import { parse } from "npm:node-html-parser@6.1.1";
import { existsSync } from "https://deno.land/std@0.158.0/node/fs.ts?s=existsSync";

interface Bookmark {
  title: string;
  url: string;
}

const bookmarks = existsSync("./data.json")
  ? JSON.parse(new TextDecoder().decode(
    Deno.readFileSync("./data.json"),
  )) as Bookmark[]
  : [] as Bookmark[];

if (bookmarks.length > 0) {
  console.log(`Loaded ${bookmarks.length} existing bookmarks`);
}

if (!existsSync("./bookmarks.html")) {
  console.log("Please place your bookmarks.html in this directory");
}

const html = new TextDecoder().decode(
  Deno.readFileSync("./bookmarks.html"),
);

const root = parse(html);

const toAdd: Bookmark[] = [];

for (const link of root.querySelectorAll("a")) {
  const { text: title, attributes } = link;
  const { href, HREF } = attributes;
  const url = href || HREF;
  if (
    title && url && title.length > 0 && url.length > 0 &&
    !bookmarks.find((el) => el.url === url)
  ) {
    toAdd.push({
      title,
      url,
    } as Bookmark);
  }
}
toAdd.reverse();

console.log(`Adding ${toAdd.length} new bookmarks`);
const data = new TextEncoder().encode(
  JSON.stringify(
    [
      ...bookmarks,
      ...toAdd,
    ],
    undefined,
    "\t",
  ),
);

Deno.writeFileSync("./data.json", data);

console.log(`Written ${data.length} bytes.`);
