"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import { Trash2 } from "lucide-react";
import React, { useState } from "react";

export default function CrawlPage() {
  const [urlInput, setUrlInput] = useState("");
  const [urlList, setUrlList] = useState<string[]>([]);
  const [isCrawling, setIsCrawling] = useState(false);
  const [crawlProgress, setCrawlProgress] = useState<number>(0);

  const addUrl = () => {
    if (urlInput.trim() && !urlList.includes(urlInput.trim())) {
      setUrlList([...urlList, urlInput.trim()]);
      setUrlInput("");
    }
  };

  const removeUrl = (url: string) => {
    setUrlList(urlList.filter((item) => item !== url));
  };

  const crawlUrls = async () => {
    if (urlList.length === 0) return;
    setIsCrawling(true);

    try {
      for (let i = 0; i < urlList.length; i++) {
        const url = urlList[i];

        // Make a request to your API to crawl the URL
        const response = await fetch("/api/crawl", {
          body: JSON.stringify({ url }),
          headers: {
            "Content-Type": "application/json",
          },
          method: "POST",
        });

        if (!response.ok) {
          console.error(`Failed to crawl ${url}`);
        }

        // Update progress
        setCrawlProgress(Math.round(((i + 1) * 100) / urlList.length));
      }

      // Handle successful crawl
      console.log("Crawling completed");
    } catch (error) {
      console.error("An error occurred during crawling:", error);
    } finally {
      setIsCrawling(false);
      setCrawlProgress(0);
    }
  };

  return (
    <div className="max-w-md mx-auto h-[80%] mt-8 stack items-center justify-center ">
      <Card
        className={"border-dashed border-2 p-6 text-center border-gray-400"}
      >
        <CardContent>
          <p className="text-gray-600">
            Add URLs to crawl and click the "Crawl" button to start crawling.
          </p>
          <div className="flex items-center space-x-2 mt-4">
            <Input
              className="flex-grow"
              onChange={(e) => setUrlInput(e.target.value)}
              placeholder="Enter URL"
              type="url"
              value={urlInput}
            />
            <Button disabled={!urlInput.trim()} onClick={addUrl}>
              Add
            </Button>
          </div>
        </CardContent>
      </Card>

      {urlList.length > 0 && (
        <div className="mt-6">
          <h3 className="text-lg font-semibold">URLs to Crawl:</h3>
          <ul className="mt-2 space-y-2">
            {urlList.map((url, idx) => (
              <li
                className="flex items-center justify-between p-2 bg-white border rounded-md"
                key={idx}
              >
                <span className="break-all">{url}</span>
                <Button
                  onClick={() => removeUrl(url)}
                  size="icon"
                  variant="ghost"
                >
                  <Trash2 className="w-5 h-5 text-red-500" />
                </Button>
              </li>
            ))}
          </ul>
          <Button
            className="mt-4 w-full"
            disabled={isCrawling}
            onClick={crawlUrls}
          >
            {isCrawling ? "Crawling..." : "Crawl"}
          </Button>
        </div>
      )}

      {isCrawling && (
        <div className="mt-6">
          <Progress value={crawlProgress} />
          <p className="text-center mt-2">{crawlProgress}%</p>
        </div>
      )}
    </div>
  );
}
