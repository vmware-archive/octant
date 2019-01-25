import moment from 'moment'
import React from 'react'

interface Props {
  config: {
    timestamp: number;
  };
  baseTime?: Date;
}

export default function Timestamp({ config, baseTime: baseDate }: Props) {
  const humanReadable = moment(config.timestamp * 1000)
    .utcOffset('+0000')
    .format('LLLL z')

  return (
    <a data-tip={humanReadable}>
      {summarizeTimestamp(config.timestamp, baseDate)}
    </a>
  )
}

/**
 * summarizeTimestamp converts a timestamp to a relative time from the current time.
 * If no date is supplied, it will use the current date.
 *
 * @param ts timestamp in seconds since epoch
 * @param base optional date to calculate from
 */
export function summarizeTimestamp(ts: number, base?: Date): string {
  let now: Date
  if (base) {
    now = base
  } else {
    now = new Date()
  }

  const then = now.getTime() / 1000 - ts

  if (then > 86400) {
    return `${Math.floor(then / 86400)}d`
  } else if (then > 3600) {
    return `${Math.floor(then / 3600)}h`
  } else if (then > 60) {
    return `${Math.floor(then / 60)}m`
  } else {
    return `${Math.floor(then)}s`
  }
}
