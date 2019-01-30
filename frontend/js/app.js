new Vue({
  el: '#app',
  data: {
    solvePage: false,
    question: '',
    answer: '',
    message: '',
    code: '',
    error: false
  },
  watch: {
    created: {
      handler() {
        const args = window.location.search.split("=", 2)
        if (args.length > 1) {
          const strings = atob(args[1]).split(':');
          this.question = strings[0];
          this.message = strings[1];

          this.solvePage = true;
        }
      },
      immediate: true,
    }
  },
  methods: {
    async createRiddle() {
      const body = {
        riddle: this.question,
        answer: this.answer,
        message: this.message
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

      this.code = btoa(`${json.riddle}:${json.message}`);
      this.error = false;
      } catch(e) {
        this.error = true;
      }
    },
    async decryptRiddle() {
      const body = {
        answer: this.answer,
        message: this.message
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

        this.code = json.message;
        this.error = false;
      } catch(e) {
        this.error = true;
      }
    },
    copy() {
      document.getElementById('share-link').select();
      document.execCommand("copy");
    }
  },
  computed: {
    shareUrl() {
      return `${window.location.href}?ridl=${this.code}`;
    }
  }
});