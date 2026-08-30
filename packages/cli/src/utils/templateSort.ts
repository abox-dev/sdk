import * as sdk from '@abox-dev/sdk'

export function sortTemplateNames<
  E extends sdk.components['schemas']['Template']['names'],
>(names: E) {
  names?.sort()
}
