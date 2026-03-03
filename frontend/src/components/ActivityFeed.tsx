import React, { useEffect, useRef } from "react";
import type { FeedItem } from "../state/types";
import ActivityFeedItem from "./ActivityFeedItem";

interface Props {
  feedItems: FeedItem[];
}

export default React.memo(function ActivityFeed({ feedItems }: Props) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "auto" });
  }, [feedItems.length]);

  return (
    <div
      ref={containerRef}
      className="flex-1 min-w-0 overflow-x-hidden overflow-y-auto p-2 space-y-1.5 h-full"
    >
      {feedItems.map((item) => (
        <ActivityFeedItem key={item.id} item={item} />
      ))}
      <div ref={bottomRef} />
    </div>
  );
});
