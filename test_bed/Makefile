.PHONY: ng react run setup

setup:
	@echo "Installing dependencies for Angular and React test beds..."
	cd angular_app && npm install
	cd react_app && npm install

ng:
	cd angular_app && npm run dev -- --port 6200

react:
	cd react_app && npm run dev -- --port 7200

run:
	@echo "Starting Angular (6200) and React (7200) apps..."
	$(MAKE) -j2 ng react

