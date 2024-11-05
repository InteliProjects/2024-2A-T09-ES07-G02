import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import '@testing-library/jest-dom';
import App from "../App";

declare var global: any;

// Mocking components
jest.mock("../components/ui/LrrContainer", () => () => <div>LrrContainer</div>);
jest.mock("../components/ui/RecentLrrs", () => ({
  RecentLrrsPopover: () => <div>RecentLrrsPopover</div>,
}));
jest.mock("../components/ui/TagsPopup", () => () => <div>TagsPopup</div>);
jest.mock("../components/ui/Filter", () => ({
  FilterPopover: () => <div>FilterPopover</div>,
}));

// Mock for mediaDevices
const mockGetUserMedia = jest.fn();
Object.defineProperty(window.navigator, "mediaDevices", {
  value: {
    getUserMedia: mockGetUserMedia,
  },
});

describe("App component", () => {
  beforeEach(() => {
    mockGetUserMedia.mockClear();
  });

  test("renders the app header and logo", () => {
    render(<App />);
    const logo = screen.getByAltText("Logo");
    const headerTitle = screen.getByText("LawHunter");

    expect(logo).toBeInTheDocument();
    expect(headerTitle).toBeInTheDocument();
  });

  test("toggles TagsPopup visibility when settings button is clicked", () => {
    render(<App />);
    
    const settingsButton = screen.getByRole("button", { name: /Settings/i }); 
  
    // Verifique se o TagsPopup não está visível inicialmente
    expect(screen.queryByText("TagsPopup")).not.toBeInTheDocument();
    
    // Simule o clique do botão
    fireEvent.click(settingsButton);
    
    // Verifique a visibilidade após o clique
    expect(screen.getByText("TagsPopup")).toBeInTheDocument(); // Ajuste conforme sua lógica
  });
  

  test("focuses search input when Ctrl+K is pressed", () => {
    render(<App />);
    const searchInput = screen.getByPlaceholderText("Digite sua pesquisa...");

    // Verify the input is not focused initially
    expect(document.activeElement).not.toBe(searchInput);

    // Simulate Ctrl+K keydown
    fireEvent.keyDown(window, { ctrlKey: true, key: "k" });

    // Verify the search input is now focused
    expect(searchInput).toHaveFocus();
  });

  test("handles Escape key to blur search input", () => {
    render(<App />);
    const searchInput = screen.getByPlaceholderText("Digite sua pesquisa...");

    // Simulate focusing the input
    searchInput.focus();
    expect(searchInput).toHaveFocus();

    // Simulate pressing the Escape key
    fireEvent.keyDown(window, { key: "Escape" });

    // Verify the input is no longer focused
    expect(document.activeElement).not.toBe(searchInput);
  });


  test("submits search when Enter key is pressed in search input", async () => {
    render(<App />);
    const searchInput = screen.getByPlaceholderText("Digite sua pesquisa...");
    const searchButton = screen.getByText("Pesquisar");

    fireEvent.change(searchInput, { target: { value: "test search" } });

    fireEvent.keyDown(searchInput, { key: "Enter" });

    // Mock the fetch request
    const mockFetch = jest.spyOn(global, "fetch").mockResolvedValue({
      ok: true,
      json: async () => ({}),
    } as Response);

    fireEvent.click(searchButton);

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith("http://localhost:8080/api/Search", expect.anything());
    });

    mockFetch.mockRestore();
  });

  test("shows RecentLrrsPopover and FilterPopover", () => {
    render(<App />);
    const recentPopover = screen.getByText("RecentLrrsPopover");
    const filterPopover = screen.getByText("FilterPopover");

    expect(recentPopover).toBeInTheDocument();
    expect(filterPopover).toBeInTheDocument();
  });

  test("toggles TagsPopup visibility when settings button is clicked", () => {
    render(<App />);
    const settingsButton = screen.getByRole("button", { name: /Settings/i });

    // Verify that the TagsPopup is not visible initially
    expect(screen.queryByText("TagsPopup")).not.toBeInTheDocument();

    // Simulate clicking the settings button
    fireEvent.click(settingsButton);

    // Verify that the TagsPopup is now visible
    expect(screen.getByText("TagsPopup")).toBeInTheDocument();

    // Simulate clicking the settings button again
    fireEvent.click(settingsButton);

    // Verify that the TagsPopup is no longer visible
    expect(screen.queryByText("TagsPopup")).not.toBeInTheDocument();
  });
});
