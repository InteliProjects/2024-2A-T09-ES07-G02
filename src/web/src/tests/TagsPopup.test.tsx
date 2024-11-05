import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import TagsPopup from '../components/ui/TagsPopup';
import '@testing-library/jest-dom/extend-expect'; // Extende expect para incluir toBeInTheDocument

describe('TagsPopup Component', () => {
  test('deve renderizar corretamente', () => {
    render(<TagsPopup onClose={() => {}} />);
    expect(screen.getByText('TAGS')).toBeInTheDocument();
  });

  test('deve abrir o input para criar nova tag', () => {
    render(<TagsPopup onClose={() => {}} />);
    fireEvent.click(screen.getByText('Criar Nova Tag'));
    expect(screen.getByPlaceholderText('Nome da nova tag')).toBeInTheDocument();
  });

  test('deve exibir mensagem quando não há tags', async () => {
    render(<TagsPopup onClose={() => {}} />);
    expect(screen.getByText('Nenhuma tag encontrada.')).toBeInTheDocument();
  });
  
  // Testar a busca
  test('deve filtrar tags pelo termo de busca', async () => {
    render(<TagsPopup onClose={() => {}} />);
    
    // Simula o termo de busca
    fireEvent.change(screen.getByPlaceholderText('Pesquisar...'), { target: { value: 'example' } });
    
    // Verifica se os resultados estão sendo filtrados corretamente
    await waitFor(() => {
      expect(screen.getByText('Nenhuma tag encontrada.')).toBeInTheDocument();
    });
  });
});
