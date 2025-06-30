'use client';

import { useEffect, useState } from 'react';

export default function Time({ rfc3339 }) {
  const [relativeTime, setRelativeTime] = useState('');

  useEffect(() => {
    const timestamp = new Date(rfc3339);

    const updateRelativeTime = () => {
      const now = new Date();
      const diffInSeconds = Math.floor((now - timestamp) / 1000);

      let display;
      let nextInterval = 60000; // default to 1 minute

      if (diffInSeconds < 60) {
        // seconds ago, update every second
        display = diffInSeconds === 1 ? '1 second ago' : `${diffInSeconds} seconds ago`;
        nextInterval = 1000; // update every second
      } else if (diffInSeconds < 3600) {
        // minutes ago, update every minute
        const minutes = Math.floor(diffInSeconds / 60);
        display = `${minutes} min${minutes !== 1 ? 's' : ''} ago`;
        nextInterval = 60000; // update every minute
      } else if (diffInSeconds < 86400) {
        // hours ago, update every hour
        const hours = Math.floor(diffInSeconds / 3600);
        display = `${hours} hour${hours !== 1 ? 's' : ''} ago`;
        nextInterval = 3600000; // update every hour
      } else {
        // older than a day, show static date, no further updates needed
        display = timestamp.toLocaleString();
        nextInterval = null; // stop updates
      }

      setRelativeTime(display);
      return nextInterval;
    };

    // Initial update
    let intervalTime = updateRelativeTime();

    if (intervalTime === null) return; // no updates needed

    const timer = setInterval(() => {
      intervalTime = updateRelativeTime();
      if (intervalTime === null) {
        clearInterval(timer);
      }
    }, intervalTime);

    // We have to clear and reset interval every time intervalTime changes,
    // so let's use a recursive timeout instead to simplify:

    clearInterval(timer);

    let timeoutId;
    const scheduleUpdate = () => {
      const next = updateRelativeTime();
      if (next !== null) {
        timeoutId = setTimeout(scheduleUpdate, next);
      }
    };
    scheduleUpdate();
    console.log(`Scheduled next update in ${intervalTime} ms`);

    return () => clearTimeout(timeoutId);
  }, [rfc3339]);

  return <time dateTime={rfc3339}>{relativeTime}</time>;
}
