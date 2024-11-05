import { Page, Text, View, Document, StyleSheet, PDFViewer } from '@react-pdf/renderer';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"; // Altere o caminho se necessário

// Estilos do PDF
const styles = StyleSheet.create({
  page: {
    flexDirection: 'column',
    backgroundColor: '#f0f0f0',
    padding: 20,
  },
  section: {
    margin: 10,
    padding: 10,
    flexGrow: 1,
    border: '1px solid #ddd',
    borderRadius: 5,
  },
  text: {
    fontSize: 14,
  },
});

const MyDocument = () => (
  <Document>
    <Page size="A4" style={styles.page}>
      <View style={styles.section}>
        <Text style={styles.text}>Seção #1: Introdução ao Visualizador de PDF</Text>
      </View>
      <View style={styles.section}>
        <Text style={styles.text}>Seção #2: Renderização Dinâmica de Conteúdo</Text>
      </View>
    </Page>
  </Document>
);

const PdfViewer = () => {
  return (
    <div style={{ width: '100%', height: '80vh', border: '1px solid #ccc', borderRadius: 8 }}>
      <PDFViewer width="100%" height="100%" style={{ border: 'none' }}>
        <MyDocument />
      </PDFViewer>
    </div>
  );
};

const PdfViewerDialog = ({ pdfUrl }: { pdfUrl: string }) => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <button
          className="flex justify-center items-center m-10 w-5/12 h-24 border border-dashed border-gray-400 rounded-2xl hover:bg-gray-100 dark:border-gray-600 dark:hover:bg-gray-700"
        >
          PDF
        </button>
      </DialogTrigger>
      <DialogContent className="w-full max-w-[800px] h-[90vh] p-4 overflow-hidden bg-gray-100 dark:bg-gray-900">
        <DialogHeader>
          <DialogTitle className="text-black dark:text-white">Visualizador de PDF</DialogTitle>
        </DialogHeader>
        <div className="h-full">
          <PdfViewer />
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default PdfViewerDialog;
