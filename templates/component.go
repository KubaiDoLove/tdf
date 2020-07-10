package templates

// ComponentTemplate represents boilerplate data for React components in TS
const ComponentTemplate = `import React, { FC } from 'react'
import { {{.Name}}Props } from './types'

const {{.Name}}: FC<{{.Name}}Props> = () => {
	return (<></>)
}

export default {{.Name}}
`

// ComponentIndexTemplate represents barrel file for React components in TS
const ComponentIndexTemplate = `export { default } from './{{.Name}}'
`

// ComponentInterfacesTemplate represents interfaces boilerplate data file for React components in TS
const ComponentInterfacesTemplate = `export interface {{.Name}}Props {}
`

// ComponentTypesIndexTemplate represents types barrel file for React components in TS
const ComponentTypesIndexTemplate = `export * from './interfaces'`
