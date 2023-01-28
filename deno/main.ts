import { writeAll } from "https://deno.land/std@0.174.0/streams/write_all.ts";

interface FeedResponseItem {
    title: string;
    copyright: string;
    fullUrl: string;
    thumbUrl: string;
    imageUrl: string;
    date: string;
    pageUrl: string;
}

async function main(): Promise<void> {
    const url = new URL("https://peapix.com/bing/feed"),
        res = await fetch(url),
        items: FeedResponseItem[] = await res.json();
    await Promise.all(
        items.map((item, index) =>
            download(new URL(item.imageUrl), index, items.length)
        )
    );
}

async function download(url: URL, index: number, total: number): Promise<void> {
    const target = url.pathname.replace("/", ""),
        res = await fetch(url),
        bytes = new Uint8Array(await res.arrayBuffer()),
        file = await Deno.create(target);
    await writeAll(file, bytes);
    file.close();
    console.log(`${index + 1}/${total}: ${target} downloaded.`);
}

await main();
