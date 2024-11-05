const fetch = require('node-fetch'); // ou a biblioteca que você estiver usando
const port = 'http://localhost:8080'; // Altere para a porta correta do seu backend

// Função para buscar tags do backend
const fetchTags = async () => {
  const response = await fetch(`${port}/api/Tag`);
  if (!response.ok) {
    throw new Error(`Erro na resposta do servidor: ${response.status}`);
  }
  return await response.json();
};

// Função para criar uma nova tag no backend
const createTag = async (tagName: string, description: string,) => {
  const response = await fetch(`${port}/api/Tag`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ tag: tagName, description }),
  });
  if (!response.ok) {
    throw new Error('Erro ao criar tag: ' + response.statusText);
  }
  return await response.json();
};

// Função para atualizar uma tag no backend
const updateTag = async (id: number, tagName: string, description: string) => {
  const response = await fetch(`${port}/api/Tag/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ id, tag: tagName, description }),
  });
  if (!response.ok) {
    throw new Error('Erro ao atualizar tag: ' + response.statusText);
  }
  return await response.json();
};

// Função para deletar uma tag no backend
const deleteTag = async (id: number) => {
  const response = await fetch(`${port}/api/Tag/${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Erro ao deletar tag: ' + response.statusText);
  }
  return id;
};

// Testes
describe('Tag API', () => {
  test('fetchTags deve retornar uma lista de tags', async () => {
    const tags = await fetchTags();
    expect(Array.isArray(tags)).toBe(true); 
    expect(tags.length).toBeGreaterThan(0); 
  });

  test('createTag deve criar uma nova tag', async () => {
    const newTag = { tag: 'Nova Tag', description: 'Descrição da nova tag' };
    const createdTag = await createTag(newTag.tag, newTag.description);
    
    expect(createdTag).toHaveProperty('id'); 
    expect(createdTag.tag).toBe(newTag.tag);
    expect(createdTag.description).toBe(newTag.description);
  });

  test('updateTag deve atualizar uma tag existente', async () => {
    const updatedTag = { id: 1, tag: 'Tag Atualizada', description: 'Descrição atualizada' };
    const result = await updateTag(updatedTag.id, updatedTag.tag, updatedTag.description);
    
    expect(result.tag).toBe(updatedTag.tag);
    expect(result.description).toBe(updatedTag.description);
  });

  test('deleteTag deve deletar uma tag existente', async () => {
    const tagId = 1;
    const result = await deleteTag(tagId);
    
    expect(result).toBe(tagId);
  });
});
