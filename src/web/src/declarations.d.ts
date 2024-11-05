// src/declarations.d.ts
declare module "*.png" {
    const value: string;
    export = value; // Alterado para usar "export =" em vez de "export default"
  }
  