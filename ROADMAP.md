<!--
 Copyright (c) 2024 Christopher Watson
 
 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# 🗺️ Roadmap

## Current Status (v0.1.0-dev)
- ✅ Core rendering engine
- ✅ Basic geometry primitives
- ✅ Color support
- ✅ Layer composition system
- ✅ Comprehensive test coverage

### Phase 1: Foundation (v0.1.0)
- [ ] Terminal Backend Implementation
  - [ ] Raw terminal mode handling
  - [ ] Input event system
  - [ ] Color and style support
  - [ ] Unicode support
  - [ ] Window resize handling

- [ ] Widget System Core
  - [ ] Widget lifecycle management
  - [ ] Layout system (Box model)
  - [ ] State management
  - [ ] Event propagation
  - [ ] Focus management

### Phase 2: Basic Widgets (v0.2.0)
- [ ] Container
- [ ] Text
- [ ] Button
- [ ] Input
- [ ] List
- [ ] Table
- [ ] Progress Bar
- [ ] Flex Layout
- [ ] Grid Layout

### Phase 3: Advanced Features (v0.3.0)
- [ ] Theming System
  - [ ] Color schemes
  - [ ] Component styles
  - [ ] Custom themes
- [ ] Animation System
  - [ ] Basic transitions
  - [ ] Keyframe animations
- [ ] Advanced Input
  - [ ] Mouse support
  - [ ] Keyboard shortcuts
  - [ ] Custom key bindings

### Phase 4: Additional Backends (v0.4.0)
- [ ] Web Backend (WASM)
  - [ ] HTML5 Canvas renderer
  - [ ] DOM event handling
- [ ] Desktop Backend
  - [ ] Cross-platform window management
  - [ ] Native input handling
- [ ] Mobile Backend (future)

### Phase 5: Developer Experience (v0.5.0)
- [ ] CLI Tools
  - [ ] Project scaffolding
  - [ ] Widget templates
  - [ ] Development server
- [ ] Hot Reload Support
- [ ] Debug Tools
  - [ ] Widget inspector
  - [ ] Performance profiler
  - [ ] State inspector

### Long-term Goals
- Custom renderer support
- Plugin system
- Accessibility features
- Rich text support
- International input methods
- Component library ecosystem

### Current Focus
We are currently focusing on Phase 1, specifically:
1. Implementing the terminal backend using tcell/termbox
2. Building the core widget system
3. Establishing the layout engine

### Contributing
Each phase has its own project board with specific tasks. Check our [Contributing Guide](CONTRIBUTING.md) for:
- How to pick up tasks
- Development setup
- Code style guidelines
- Pull request process

We welcome contributions at any level:
- 🐛 Bug fixes
- 📝 Documentation
- ✨ New features
- 🎨 Design improvements
- 💡 Suggestions and ideas

Track our progress and contribute on our [GitHub Projects page](https://github.com/watzon/tide/projects).