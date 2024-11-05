import React from "react";
import LrrCard from "./LrrCard";
import nosearch from "../../assets/Search-bro.svg";

interface LrrContainerProps {
  searchTerm: string;
  isLoading: boolean;
  searchResults: Array<{ 
    rule: string; 
    agency: string; 
    s3_url: string; 
    tags: string; 
  }>;
}

const LrrContainer: React.FC<LrrContainerProps> = ({ searchTerm, isLoading, searchResults }) => {
  console.log({searchResults})

  return (
    <div className="flex flex-col justify-center items-center w-6/12 mx-auto gap-4">
      {isLoading ? (
        <div className="flex flex-col justify-center items-center mt-48 h-max w-max">
          <p className="text-3xl text-gray-700 dark:text-gray-200 font-semibold mb-8">Pesquisando...</p>
          <div className="spinner-border animate-spin inline-block w-20 h-20 border-4 rounded-full mt-8"></div>
        </div>
      ) : searchTerm && searchResults.length > 0 ? (
        searchResults.map((lrr, index) => {
          console.log({lrr})
          return (
          <LrrCard
            key={index}
            rule={lrr.rule}
            agency={lrr.agency}
            s3_url={lrr.s3_url}
            tags={[lrr.tags]}
          />
        )})
      ) : (
        <div className="flex flex-col justify-center items-center">
          <img src={nosearch} alt="Nenhum resultado encontrado" />
          <p className="text-lg text-gray-700 dark:text-gray-200 font-semibold">
            Parece que você ainda não pesquisou nada.
          </p>
          <p className="text-md text-gray-500 dark:text-gray-400">
            Digite algo na barra de pesquisa acima para encontrar LRRs.
          </p>
        </div>
      )}
    </div>
  );
};

export default LrrContainer;
