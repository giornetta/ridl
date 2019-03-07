new Vue({
  el: '#app',
  data: {
    solvePage: false,
    question: '',
    answer: '',
    message: '',
    ignoreCase: false,
    ignoreSpaces: false,
    id: '',
    error: false
  },
  watch: {
    created: {
      async handler() {
        const args = window.location.search.split("=", 2)
        if (args.length > 1) {
          try {
            const res = await fetch(`http://localhost:8080/ridl/riddle/${args[1]}`);
            const json = await res.json();
            
            this.question = json.question;
            this.id = args[1];
            this.error = false;
            this.solvePage = true;
            } catch(e) {
              this.showError();
            }
        }
      },
      immediate: true,
    }
  },
  methods: {
    async createRiddle() {
      const body = {
        question: this.question,
        answer: this.answer,
        message: this.message,
        ignoreCase: this.ignoreCase,
        ignoreSpaces: this.ignoreSpaces
      };
      
      try {
      const res = await fetch('http://localhost:8080/ridl/encrypt', {
        headers: {
          'Content-Type': 'application/json'
        },
        method: 'POST',
        body: JSON.stringify(body)
      });

      const json = await res.json();

      this.id = json.riddleID;
      this.error = false;
      } catch(e) {
        this.showError();
      }
    },
    async decryptRiddle() {
      const body = {
        riddleID: this.id,
        answer: this.answer,
      };
    
      try {
        const res = await fetch('http://localhost:8080/ridl/decrypt', {
          headers: {
            'Content-Type': 'application/json'
          },
          method: 'POST',
          body: JSON.stringify(body)
        });
        
        if (res.status !== 200) {
          throw res.status;
        }

        const json = await res.json();
        this.message = json.message;
        console.log(this.message);
        this.error = false;
      } catch(e) {
        this.showError();
      }
    },
    copy() {
      document.getElementById('share-link').select();
      document.execCommand("copy");
    },
    showError() {
      this.error = true;
      setTimeout(() => {
        this.error = false;
      }, 2000);
    }
  },
  computed: {
    shareUrl() {
      return `${window.location.href}?ridl=${this.id}`;
    }
  }
});