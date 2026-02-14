import antfu from '@antfu/eslint-config';
import prettier from 'eslint-config-prettier';

export default antfu(
    {
        stylistic: false,
        rules: {
            'antfu/no-top-level-await': 'off',
            'ts/method-signature-style': 'off',
            'ts/no-empty-object-type': 'off',
            'ts/no-namespace': 'off',
        },
    },
    prettier,
);
