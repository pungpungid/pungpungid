// Main application JavaScript
document.addEventListener('DOMContentLoaded', function() {
  // Sidebar toggle for mobile
  const sidebarToggle = document.getElementById('sidebar-toggle');
  const sidebar = document.getElementById('sidebar');
  
  if (sidebarToggle && sidebar) {
    sidebarToggle.addEventListener('click', function() {
      sidebar.classList.toggle('open');
    });
  }

  // Close sidebar when clicking outside on mobile
  document.addEventListener('click', function(event) {
    if (window.innerWidth < 768) {
      if (!sidebar.contains(event.target) && !sidebarToggle.contains(event.target)) {
        sidebar.classList.remove('open');
      }
    }
  });

  // Auto-hide alerts
  const alerts = document.querySelectorAll('.alert');
  alerts.forEach(function(alert) {
    setTimeout(function() {
      alert.style.opacity = '0';
      alert.style.transform = 'translateY(-10px)';
      setTimeout(function() {
        alert.remove();
      }, 300);
    }, 5000);
  });

  // Form submission with AJAX
  const forms = document.querySelectorAll('.ajax-form');
  forms.forEach(function(form) {
    form.addEventListener('submit', function(e) {
      e.preventDefault();
      
      const submitBtn = form.querySelector('button[type="submit"]');
      const originalText = submitBtn.innerHTML;
      submitBtn.disabled = true;
      submitBtn.innerHTML = '<span class="spinner"></span> Loading...';
      
      fetch(form.action, {
        method: form.method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(Object.fromEntries(new FormData(form))),
      })
      .then(response => response.json())
      .then(data => {
        if (data.redirect) {
          window.location.href = data.redirect;
        } else if (data.error) {
          showToast(data.error, 'error');
        } else {
          showToast('Success!', 'success');
        }
      })
      .catch(error => {
        showToast('An error occurred', 'error');
        console.error('Error:', error);
      })
      .finally(() => {
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;
      });
    });
  });

  // Confirm delete actions
  const deleteButtons = document.querySelectorAll('.btn-delete');
  deleteButtons.forEach(function(btn) {
    btn.addEventListener('click', function(e) {
      if (!confirm('Are you sure you want to delete this item?')) {
        e.preventDefault();
      }
    });
  });
});

// Toast notification function
function showToast(message, type = 'info') {
  const toast = document.createElement('div');
  toast.className = `toast ${type === 'error' ? 'bg-red-500 text-white' : type === 'success' ? 'bg-green-500 text-white' : 'bg-blue-500 text-white'}`;
  toast.textContent = message;
  document.body.appendChild(toast);
  
  setTimeout(function() {
    toast.style.opacity = '0';
    toast.style.transform = 'translateY(10px)';
    setTimeout(function() {
      toast.remove();
    }, 300);
  }, 3000);
}

// Loading state helper
function setLoading(element, loading) {
  if (loading) {
    element.setAttribute('disabled', 'true');
    element.classList.add('opacity-50', 'cursor-not-allowed');
  } else {
    element.removeAttribute('disabled');
    element.classList.remove('opacity-50', 'cursor-not-allowed');
  }
}
