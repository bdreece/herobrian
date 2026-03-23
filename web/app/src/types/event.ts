declare global {
    namespace Herobrian {
        interface Event {
            type: string;
            data: unknown;
        }
    }
}

export {};
