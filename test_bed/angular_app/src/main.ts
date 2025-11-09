import '@angular/compiler'
import 'zone.js'
import { bootstrapApplication } from '@angular/platform-browser'
import { AppComponent } from './app.component'

import '@rx/micro-prescription-info'
import '@rx/micro-patient-info'

bootstrapApplication(AppComponent)
  .then(() => console.log('Angular test bed bootstrapped'))
  .catch(err => console.error(err))

