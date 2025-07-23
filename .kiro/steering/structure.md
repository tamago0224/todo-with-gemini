# Project Structure Steering: Todo App

## 1. Directory Organization

```
/
├── public/         # Static assets
├── src/
│   ├── components/   # Reusable UI components
│   ├── pages/        # Page components
│   ├── services/     # API interaction logic
│   ├── styles/       # Global styles and themes
│   ├── utils/        # Utility functions
│   └── App.tsx       # Main application component
├── .eslintrc.js    # ESLint configuration
├── .gitignore      # Git ignore file
├── package.json    # Project dependencies and scripts
├── tsconfig.json   # TypeScript configuration
└── README.md       # Project documentation
```

## 2. Naming Conventions

- **Components**: PascalCase (e.g., `TodoItem.tsx`)
- **Files/Folders**: kebab-case (e.g., `todo-list`)
- **Variables/Functions**: camelCase (e.g., `addTask`)

## 3. Code Style

- Follow standard TypeScript and React best practices.
- Use ESLint and Prettier for code linting and formatting.
