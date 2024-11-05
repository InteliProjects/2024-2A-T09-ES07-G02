from gensim.models import KeyedVectors
import numpy as np

class WordVectorizer:
    def __init__(self, model_path):
        """
        Initializes the WordVectorizer by loading the Word2Vec model.
        
        :param model_path: Path to the Word2Vec model file.
        """
        self.model = KeyedVectors.load_word2vec_format(model_path, binary=False)
        
    def vectorize_token(self, token):
        """
        Returns the vector for a single token if it is in the vocabulary.
        
        :param token: Word or token to be vectorized.
        :return: Word2Vec vector for the token, or None if not in vocabulary.
        """
        if token in self.model:
            return self.model[token]
        else:
            return None

    def average_vector(self, tokens):
        """
        Computes the average vector for a list of tokens, ignoring those not in the vocabulary.
        
        :param tokens: List of words (tokens) to compute the average vector.
        :return: Average vector of the list of tokens.
        """
        # Vectorize only tokens that are in the model's vocabulary
        valid_vectors = [self.vectorize_token(token) for token in tokens if self.vectorize_token(token) is not None]
        
        if valid_vectors:
            # Compute the mean of the vectors
            mean_vector = np.mean(valid_vectors, axis=0)
            return mean_vector.reshape(1, -1)
        else:
            return np.zeros(self.model.vector_size)  # Returns a zero vector if no words are found
