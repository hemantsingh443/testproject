const resultsDiv = document.getElementById('results');
const searchForm = document.getElementById('search-form');

function createPaperEntry(paper) {
  const listItem = document.createElement('li');

  // View on arXiv Link
  const paperLink = document.createElement('a');
  paperLink.href = paper.pdf_url.replace('.pdf', '');
  paperLink.textContent = 'View on arXiv';
  paperLink.target = '_blank';
  listItem.appendChild(paperLink);

  // Paper details
  listItem.innerHTML += `
    <h3>${paper.title}</h3>
    <p>Authors: ${paper.authors ? paper.authors.join(', ') : 'N/A'}</p>
    <p>${paper.summary}</p>
  `;

  // Summarize Button
  const summarizeButton = document.createElement('button');
  summarizeButton.textContent = 'Summarize';
  summarizeButton.classList.add('summarize-button');
  summarizeButton.dataset.pdfUrl = paper.pdf_url;
  listItem.appendChild(summarizeButton);

  // Summary Box
  const summaryBox = document.createElement('div');
  summaryBox.classList.add('summary-box');
  listItem.appendChild(summaryBox);

  resultsDiv.appendChild(listItem);
}

function performSearch(event) {
  event.preventDefault();
  const searchQuery = document.getElementById('search_query').value;

  fetch('/search', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ search_query: searchQuery })
  })
  .then(response => response.json())
  .then(data => {
    resultsDiv.innerHTML = ''; // Clear previous results

    const papers = data.results;
    if (papers && papers.length > 0) {
        papers.forEach(paper => {
            createPaperEntry(paper);
        });
    } else {
        resultsDiv.textContent = 'No papers found.'
    }
  })
  .catch(error => console.error('Error:', error));
}

document.addEventListener('click', function(event) {
  if (event.target.classList.contains('summarize-button')) {
    const button = event.target;
    const pdfUrl = button.dataset.pdfUrl;
    const summaryBox = button.nextElementSibling;

    fetch('/summarize', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ pdf_url: pdfUrl })
    })
    .then(response => response.text())
    .then(summary => {
      summaryBox.textContent = summary;
      summaryBox.style.display = 'block';
    })
    .catch(error => console.error('Error:', error));
  }
});

searchForm.addEventListener('submit', performSearch);