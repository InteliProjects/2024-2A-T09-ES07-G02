// src/@types/jest-fetch-mock.d.ts
declare module 'jest-fetch-mock' {
    function enableMocks(): void;
    function enableFetchMocks(): void;
    function disableMocks(): void;
  
    const fetchMock: {
      mockImplementation: typeof jest.fn;
      mockImplementationOnce: typeof jest.fn;
      resetMocks: () => void;
      // Adicione outros métodos que você deseja usar
    };
  
    export { enableMocks, enableFetchMocks, disableMocks, fetchMock };
  }
  