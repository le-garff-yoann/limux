import { mount } from '@vue/test-utils' 
import App from '@/App.vue'

describe('App.vue', () => {
  it('has the toYaml filters', () => {
    expect(typeof App.filters.toYaml).toBe('function')
  })

  it('has a data hook', () => {
    expect(typeof App.data).toBe('function')
  })

  it('has a created hook', () => {
    expect(typeof App.created).toBe('function')
  })

  const data = App.data()

  it('sets the correct default data', () => {
    expect(typeof data.processors).toBe('object')
    expect(typeof data.ws).toBe('object')
    expect(typeof data.ws.isConnected).toBe('boolean')
    expect(data.ws.isConnected).toBeFalsy()
  })

  it('has the pushEvent method', () => {
    expect(typeof App.methods.pushEvent).toBe('function')
  })

  const wrapper = mount(App)

  it('is a Vue instance', () => {
    expect(wrapper.isVueInstance()).toBeTruthy()
  })

  const app = wrapper.find(App)

  it('renders app', () => {
    expect(app.is('div')).toBe(true)
  })

  it('renders processors', () => {
    const processors = [
      { processor: { foo: "foo" }, type: 0 },
      { processor: { bar: "bar" }, type: 0 }
    ]

    for (const pr of processors) {
      wrapper.vm.pushEvent(pr)
    }

    expect(app.findAll('.list-group-item').length).toBe(processors.length)
  })
})
