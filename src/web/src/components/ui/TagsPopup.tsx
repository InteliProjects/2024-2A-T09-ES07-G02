import React, { useState, useEffect } from 'react';

export interface Tag {
  id: number;
  tag: string;
  description: string;
}

interface TagsPopupProps {
  onClose: () => void;
}

const TagsPopup: React.FC<TagsPopupProps> = ({ onClose }) => {
  const [tags, setTags] = useState<Tag[]>([]);
  const [editingTag, setEditingTag] = useState<number | null>(null);
  const [editValue, setEditValue] = useState(''); // State to store the tag's temporary value
  const [editDescription, setEditDescription] = useState(''); // State to store the tag's temporary description
  const [newTagName, setNewTagName] = useState<string>(''); // State for the new tag name
  const [newTagDescription, setNewTagDescription] = useState<string>(''); // State for the new tag description
  const [isAddingNewTag, setIsAddingNewTag] = useState(false); // State to control the new tag input
  const [searchTerm, setSearchTerm] = useState(''); // State to store the search term

  const port = `http://localhost:8080`;

  // Function to fetch tags from the backend
  const fetchTags = async (): Promise<Tag[]> => {
    try {
      const response = await fetch(`${port}/api/Tag`);
      if (!response.ok) {
        throw new Error(`Server response error: ${response.status}`);
      }
      const data: Tag[] = await response.json();
      setTags(data);
      return data;
    } catch (error) {
      console.error('Error fetching tags:', error);
      return [];
    }
  };

  // Load tags from the backend when the component mounts
  useEffect(() => {
    fetchTags();
  }, []);

  // Function to create a new tag in the backend
  const createTag = async (tagName: string, description: string): Promise<void> => {
    try {
      const response = await fetch(`${port}/api/Tag`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ tag: tagName, description }),
      });
      if (!response.ok) {
        throw new Error('Error creating tag: ' + response.statusText);
      }
      fetchTags();
      setIsAddingNewTag(false);
    } catch (error) {
      console.error('Error creating tag:', error);
    }
  };

  // Function to update a tag in the backend
  const updateTag = async (id: number, tagName: string, description: string) => {
    try {
      const response = await fetch(`${port}/api/Tag/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id, tag: tagName, description }),
      });
      if (!response.ok) {
        throw new Error('Error updating tag: ' + response.statusText);
      }
      fetchTags();
    } catch (error) {
      console.error('Error updating tag:', error);
    }
  };

  // Function to delete a tag from the backend
  const deleteTag = async (id: number) => {
    try {
      const response = await fetch(`${port}/api/Tag/${id}`, {
        method: 'DELETE',
      });
      if (!response.ok) {
        throw new Error('Error deleting tag: ' + response.statusText);
      }
      fetchTags();
    } catch (error) {
      console.error('Error deleting tag:', error);
    }
  };

  // Start editing a tag
  const startEditingTag = (index: number) => {
    setEditingTag(index);
    setEditValue(tags[index].tag);
    setEditDescription(tags[index].description);
  };

  // Save tag on Enter key press
  const saveTagOnEnter = (event: React.KeyboardEvent<HTMLInputElement>, index: number) => {
    if (event.key === 'Enter') {
      const tagId = tags[index].id;
      updateTag(tagId, editValue, editDescription);
      setEditingTag(null);
    }
  };

  // Delete a tag
  const handleDeleteTag = (index: number) => {
    const tagId = tags[index].id;
    deleteTag(tagId);
  };

  // Add a new tag
  const handleAddNewTag = () => {
    if (newTagName.trim() !== '' && newTagDescription.trim() !== '') {
      createTag(newTagName.trim(), newTagDescription.trim());
      setNewTagName('');
      setNewTagDescription('');
    }
  };

  const saveNewTagOnEnter = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter') {
      handleAddNewTag();
    }
  };

  // Filter tags based on the search term (tag name)
  const filteredTags = tags.filter((tag) =>
    tag && tag.tag && tag.tag.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
      <div className="p-6 rounded-lg relative w-full max-w-md border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800 text-black dark:text-white max-h-[80vh] mt-8 overflow-auto">
        {/* TAGS and Close Button */}
        <div className="flex justify-between items-center mb-4 border-b pb-2 border-gray-300 dark:border-gray-700">
          <span className="text-lg font-bold">TAGS</span>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-300">
            ✕
          </button>
        </div>

        {/* Add Tag and Search */}
        <div className="flex justify-between items-center mb-4">
          <button
            onClick={() => setIsAddingNewTag(true)}
            className="bg-blue-500 text-white hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700 px-4 py-2 rounded"
          >
            Create New Tag
          </button>
          <div className="relative">
            <input
              type="text"
              placeholder="Search..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="border rounded pl-8 pr-4 py-1 bg-gray-100 dark:bg-gray-700 text-black dark:text-white focus:outline-none"
            />
            <svg
              className="absolute left-2 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M21 21l-4.35-4.35m0 0a7.5 7.5 0 111.45-1.45L21 21z"
              ></path>
            </svg>
          </div>
        </div>

        {/* Tags Table */}
        <div className="flex flex-col space-y-2 overflow-auto max-h-[40vh]">
          {filteredTags.length > 0 ? (
            filteredTags.map((tag, index) => (
              <div key={tag.id} className="flex justify-between items-center border p-2 rounded-md border-gray-300 dark:border-gray-700">
                {editingTag === index ? (
                  <>
                    <input
                      value={editValue}
                      onChange={(e) => setEditValue(e.target.value)}
                      onKeyDown={(e) => saveTagOnEnter(e, index)}
                      className="flex-grow border rounded py-1 px-2 bg-gray-100 dark:bg-gray-700 text-black dark:text-white focus:outline-none"
                      autoFocus
                    />
                    <input
                      value={editDescription}
                      onChange={(e) => setEditDescription(e.target.value)}
                      onKeyDown={(e) => saveTagOnEnter(e, index)}
                      className="flex-grow border rounded py-1 px-2 bg-gray-100 dark:bg-gray-700 text-black dark:text-white focus:outline-none"
                      placeholder="Description"
                    />
                    <button onClick={() => setEditingTag(null)} className="text-red-400 hover:text-red-600">
                      Cancel
                    </button>
                  </>
                ) : (
                  <>
                    <div className="flex-grow">
                      <strong>{tag.tag}</strong>
                      <p className="text-gray-400 text-sm">{tag.description}</p>
                    </div>
                    <div className="flex space-x-2">
                      <button onClick={() => startEditingTag(index)} className="text-blue-400 hover:text-blue-600">
                        Edit
                      </button>
                      <button onClick={() => handleDeleteTag(index)} className="text-red-400 hover:text-red-600">
                        Delete
                      </button>
                    </div>
                  </>
                )}
              </div>
            ))
          ) : (
            <div>No tags found.</div>
          )}
        </div>

        {/* Add New Tag Form */}
        {isAddingNewTag && (
          <div className="flex flex-col mt-4 space-y-2">
            <input
              value={newTagName}
              onChange={(e) => setNewTagName(e.target.value)}
              placeholder="New Tag Name"
              onKeyDown={saveNewTagOnEnter}
              className="border rounded py-2 px-4 bg-gray-100 dark:bg-gray-700 text-black dark:text-white focus:outline-none"
            />
            <input
              value={newTagDescription}
              onChange={(e) => setNewTagDescription(e.target.value)}
              placeholder="New Tag Description"
              onKeyDown={saveNewTagOnEnter}
              className="border rounded py-2 px-4 bg-gray-100 dark:bg-gray-700 text-black dark:text-white focus:outline-none"
            />
            <button
              onClick={handleAddNewTag}
              className="bg-blue-500 text-white hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700 px-4 py-2 rounded"
            >
              Add Tag
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default TagsPopup;
