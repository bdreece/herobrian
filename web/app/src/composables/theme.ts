import {
    type BasicColorSchema,
    makeDestructurable,
    useColorMode,
    type UseColorModeReturn,
    useToggle,
} from '@vueuse/core';

export type UseThemeReturn = {
    readonly theme: UseColorModeReturn<'light' | 'dark'>;
    readonly toggle: (value?: BasicColorSchema | undefined) => BasicColorSchema;
} & readonly [
    UseColorModeReturn<'light' | 'dark'>,
    (value?: BasicColorSchema | undefined) => BasicColorSchema,
];

export function useTheme(): UseThemeReturn {
    const theme = useColorMode({
        storageKey: 'herobrian_theme',
        attribute: 'data-theme',
        modes: {
            light: 'autumn',
            dark: 'coffee',
        },
    });

    const toggle = useToggle(theme, {
        truthyValue: 'light',
        falsyValue: 'dark',
    });

    return makeDestructurable(
        { theme, toggle } as const,
        [theme, toggle] as const,
    );
}
