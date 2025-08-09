import os
import base64

def pdfs_to_base64_txt(root_dir, output_file):
    pdf_files = [f for f in os.listdir(root_dir) if f.lower().endswith('.pdf')]
    with open(output_file, 'w', encoding='utf-8') as out_f:
        for pdf_file in pdf_files:
            pdf_path = os.path.join(root_dir, pdf_file)
            with open(pdf_path, 'rb') as f:
                encoded = base64.b64encode(f.read()).decode('utf-8')
                out_f.write(encoded + '\n\n')

if __name__ == '__main__':
    root_directory = '.'
    output_txt = 'pdf_base64.txt'
    pdfs_to_base64_txt(root_directory, output_txt)