import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import Apply from '@/pages/Apply.vue'

describe('Apply.vue', () => {
  it('renders props.name when passed', () => {
    const wrapper = mount(Apply, {  })
    expect(wrapper.text()).toContain('NePlus')
  })
})
