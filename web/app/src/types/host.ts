export interface Host {
    arch: string;
    platform: string;
    image: {
        type: string;
        memory: number;
        network: string;
        storage: null;
    };
    processor: {
        cores: number;
        threads: number;
    };
}
