import React from "react";
import type { FeedItem } from "../state/types";
import ActivityFeedItem from "./ActivityFeedItem";

interface Props {
  feedItems: FeedItem[];
}

export default React.memo(function ActivityFeed({ feedItems }: Props) {
  return (
    <div className="min-w-0 space-y-1.5">
      {feedItems.map((item) => (
        <ActivityFeedItem key={item.id} item={item} />
      ))}
    </div>
  );
});
