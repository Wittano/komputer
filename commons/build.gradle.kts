import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.9.10"
    // TODO Replace kapt by ksp
    kotlin("kapt") version "1.9.10"
}

group = "com.wittano.komputer"
version = rootProject.version

repositories {
    mavenCentral()
}

val jacksonVersion = "2.15.2"
val daggerVersion = "2.48"

dependencies {
    // Utils
    implementation("io.github.cdimascio:dotenv-kotlin:6.4.1")

    testImplementation(platform("org.junit:junit-bom:5.9.1"))
    testImplementation("org.junit.jupiter:junit-jupiter")
}

tasks.test {
    useJUnitPlatform()
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "17"
}