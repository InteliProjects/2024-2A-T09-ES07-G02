// global.d.ts
export {};

declare global {
    var sendToBackend: (audioBlob: Blob, transcript: string) => Promise<void>;
}
