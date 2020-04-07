export default {
  initialColorMode: 'light',
  colors: {
    text: 'rgb(70, 78, 78)',
    background: '#fff',
    primary: '#007E8A',
    secondary: '#007E8A',
    muted: '#f6f6f6',
    highlight: '#ffffcc',
    gray: '#777',
    purple: '#609',
    heading: '#000d0d',
    modes: {
      dark: {
        text: '#D8DEE9',
        background: '#2E3440',
        primary: '#70ccd2',
        secondary: '#70ccd2',
        muted: '#c2d6d82e',
        highlight: '#ffffcc',
        gray: '#999',
        purple: '#c0f',
        heading: '#fafbfb',
      }
    },
  },
  fonts: {
    body: 'Open Sans',
    heading: 'inherit',
    monospace: 'Menlo, monospace',
  },
  fontSizes: [12, 14, 16, 20, 24, 32, 48, 64, 72],
  fontWeights: {
    body: '400',
    heading: '600',
  },
  lineHeights: {
    body: 1.65,
    heading: 1.25,
  },
  textStyles: {
    heading: {
      fontFamily: 'heading',
      fontWeight: 'heading',
      lineHeight: 'heading',
      a: {
        color: 'inherit',
        textDecoration: 'none'
      }
    },
    display: {
      variant: 'textStyles.heading',
      fontSize: [5],
      mt: 3,
    },
  },
  styles: {
    Container: {
      p: 3,
      maxWidth: 1024,
    },
    root: {
      fontFamily: 'body',
      lineHeight: 'body',
      fontWeight: 'body',
      fontSize: 1
    },
    h1: {
      variant: 'textStyles.display',
      color: 'heading'
    },
    h2: {
      variant: 'textStyles.heading',
      fontSize: 4,
      color: 'text'
    },
    h3: {
      variant: 'textStyles.heading',
      fontSize: 3,
      color: 'text'
    },
    h4: {
      variant: 'textStyles.heading',
      fontSize: 2,
      color: 'text'
    },
    h5: {
      variant: 'textStyles.heading',
      fontSize: 1,
    },
    h6: {
      variant: 'textStyles.heading',
      fontSize: 0,
    },
    a: {
      color: 'primary',
      '&:hover': {
        color: 'secondary',
      },
    },
    pre: {
      variant: 'prism',
      fontFamily: 'monospace',
      fontSize: 1,
      p: 3,
      color: 'text',
      bg: 'muted',
      overflow: 'auto',
      code: {
        color: 'inherit',
      },
    },
    code: {
      fontFamily: 'monospace',
      color: 'primary',
      BackgroundColor: 'muted',
      fontSize: 1,
      padding: '30px'
    },
    inlineCode: {
      fontFamily: 'monospace',
      color: 'secondary',
      bg: 'muted',
    },
    table: {
      width: '100%',
      my: 4,
      borderCollapse: 'separate',
      borderSpacing: 0,
      [['th', 'td']]: {
        textAlign: 'left',
        py: '4px',
        pr: '4px',
        pl: 0,
        borderColor: 'muted',
        borderBottomStyle: 'solid',
      },
    },
    th: {
      verticalAlign: 'bottom',
      borderBottomWidth: '2px',
    },
    td: {
      verticalAlign: 'top',
      borderBottomWidth: '1px',
    },
    hr: {
      border: 0,
      borderBottom: '1px solid',
      borderColor: 'muted',
    },
    img: {
      maxWidth: '100%'
    },
    ul: {
      fonts: 'body',
      color: 'text',
    },
    p: {
      fonts: 'body',
      color: 'text',
    }
  },
  prism: {
    [[
      '.comment',
      '.prolog',
      '.doctype',
      '.cdata',
      '.punctuation',
      '.operator',
      '.entity',
      '.url',
    ]]: {
      color: 'gray',
    },
    '.comment': {
      fontStyle: 'italic',
    },
    [[
      '.property',
      '.tag',
      '.boolean',
      '.number',
      '.constant',
      '.symbol',
      '.deleted',
      '.function',
      '.class-name',
      '.regex',
      '.important',
      '.variable',
    ]]: {
      color: 'purple',
    },
    [['.atrule', '.attr-value', '.keyword']]: {
      color: 'primary',
    },
    [[
      '.selector',
      '.attr-name',
      '.string',
      '.char',
      '.builtin',
      '.inserted',
    ]]: {
      color: 'secondary',
    },
  },
}
