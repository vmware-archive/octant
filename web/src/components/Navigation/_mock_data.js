const NavigationData = {
  sections: [
    {
      name: 'Overview',
      link: '/'
    },
    {
      name: 'Section 1',
      link: '#section1',
      children: [
        {
          name: 'Item 1',
          link: '#item1'
        },
        {
          name: 'Item 2',
          link: '#item2'
        },
        {
          name: 'Item 3',
          link: '#item3'
        },
        {
          name: 'Item 4',
          link: '#item4'
        }
      ]
    },
    {
      name: 'Section2',
      link: '#section2'
    }
  ]
}

export default NavigationData
