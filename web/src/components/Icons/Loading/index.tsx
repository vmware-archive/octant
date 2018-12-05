import React from 'react'
import './styles.scss'

interface Props {
  className?: string;
}

export default function({ className = 'loading' }: Props) {
  return (
    <svg
      className={className}
      xmlns='http://www.w3.org/2000/svg'
      width='93'
      height='89'
    >
      <g fill='#FFF' fillRule='nonzero' stroke='#F2582D' strokeWidth='3'>
        <path
          d='M42 15L66.058132 26.4876113 72 52.3000396 55.351256 73 28.648744 73 12 52.3000396 17.941868 26.4876113z'
          transform='translate(3 3)'
        />
        <path
          d={'M43.5 0L78.3842914 16.6372302 87 54.0207471 62.8593212 84 ' +
          '24.1406788 84 0 54.0207471 8.61570857 16.6372302z'}
          transform='translate(3 3)'
        />
      </g>
    </svg>
  )
}
