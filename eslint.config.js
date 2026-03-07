import antfu from '@antfu/eslint-config';
import prettier from 'eslint-config-prettier';

export default antfu(
    {
        ignores: ['**/*.d.ts', '**/dist', '**/coverage'],
        stylistic: false,
        rules: {
            'dot-notation': 'off',
            'antfu/no-top-level-await': 'off',
            'import/consistent-type-specifier-style': 'off',
            'ts/method-signature-style': 'off',
            'ts/no-empty-object-type': 'off',
            'ts/no-namespace': 'off',
            'vue/no-ref-as-operand': 'off',
        },
    },
    prettier,
);
