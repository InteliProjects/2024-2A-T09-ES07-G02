import { Filter } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Badge } from "./badge";
import { useState } from "react";

const availableFilters = [
  { name: 'Órgão Regulador', options: ['ANBIMA', 'Banco Central', 'CVM', 'Outros'] },
  { name: 'Status', options: ['Ativo', 'Revogado', 'Pendente'] },
];

const filterColors: Record<string, string> = {
  "ANBIMA": "bg-yellow-100 text-yellow-800 dark:bg-yellow-600 dark:text-yellow-100",
  "Banco Central": "bg-purple-100 text-purple-800 dark:bg-purple-600 dark:text-purple-100",
  "CVM": "bg-teal-100 text-teal-800 dark:bg-teal-600 dark:text-teal-100",
  "Outros": "bg-pink-100 text-pink-800 dark:bg-pink-600 dark:text-pink-100",
  "Ativo": "bg-green-200 text-green-900 dark:bg-green-600 dark:text-green-100",
  "Revogado": "bg-red-200 text-red-900 dark:bg-red-600 dark:text-red-100",
  "Pendente": "bg-orange-200 text-orange-900 dark:bg-orange-600 dark:text-orange-100",
};

export function FilterPopover() {
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');

  const [selectedFilters, setSelectedFilters] = useState({
    regulator: [],
    status: [],
  });

  const handleSubmit = () => {
    const filters = { startDate, endDate, selectedFilters };
    console.log("Filtros enviados:", filters);
  };

  return (
    <Popover>
      <PopoverTrigger asChild>
        <button className="relative p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 focus:outline-none">
          <Filter className="h-6 w-6 text-black dark:text-white" />
        </button>
      </PopoverTrigger>
      <PopoverContent className="w-[28rem] p-4 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm bg-white dark:bg-gray-800">
        <h4 className="font-semibold mb-2 border-b pb-2 text-gray-700 dark:text-gray-200">Filtros Disponíveis</h4>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Data de Início</label>
          <input
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            className="w-full p-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-100"
          />
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Data Final</label>
          <input
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
            className="w-full p-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-100"
          />
        </div>

        <div className="max-h-72 overflow-y-auto mb-4">
          <ul className="space-y-3">
            {availableFilters.map((filter, index) => (
              <li
                key={index}
                className="flex flex-col p-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 shadow-sm"
              >
                <span className="text-lg font-medium mb-2 text-gray-700 dark:text-gray-200">{filter.name}</span>
                <div className="flex flex-wrap gap-2">
                  {filter.options.map((option, idx) => (
                    <Badge
                      key={idx}
                      className={`${filterColors[option] || "bg-gray-100 text-gray-800 dark:bg-gray-600 dark:text-gray-300"} whitespace-nowrap`}
                    >
                      {option}
                    </Badge>
                  ))}
                </div>
              </li>
            ))}
          </ul>
        </div>

        <div className="flex justify-center">
          <button
            onClick={handleSubmit}
            className="bg-blue-500 dark:bg-blue-700 hover:bg-blue-700 dark:hover:bg-blue-900 text-white font-semibold py-2 px-6 rounded-full shadow-md transition-all duration-300"
          >
            Filtrar
          </button>
        </div>
      </PopoverContent>
    </Popover>
  );
}
