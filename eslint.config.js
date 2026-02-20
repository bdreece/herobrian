import antfu from '@antfu/eslint-config';
import prettier from 'eslint-config-prettier';

export default antfu(
    {
        ignores: ['**/*.d.ts', '**/dist', '**/coverage'],
        stylistic: false,
        rules: {
            'antfu/no-top-level-await': 'off',
            'ts/method-signature-style': 'off',
            'ts/no-empty-object-type': 'off',
            'ts/no-namespace': 'off',
            'import/consistent-type-specifier-style': 'off',
        },
    },
    prettier,
);
