import { Template } from '@abox-dev/sdk'

export const template = Template()
  .fromImage('alpine:latest')
  .setUser('root')
  .setWorkdir('/')
  .copy('package.json', '/app/')
  .copy('src/index.js', './src/')
  .copy('config.json', '/etc/app/config.json')
  .setUser('user')
  .setWorkdir('/home/user')
